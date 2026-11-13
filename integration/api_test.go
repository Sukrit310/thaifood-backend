// ไฟล์ integration/api_test.go: ทดสอบครบเส้นทาง HTTP → middleware → handler → PostgreSQL ด้วย server และ schema ชั่วคราว
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"thaifood-backend/auth"
	"thaifood-backend/handler"
	"thaifood-backend/mw"
	"thaifood-backend/store"
)

// secret นี้ใช้เฉพาะ JWT ของ server ทดสอบที่สร้างภายใน process นี้
const testSecret = "integration-only-secret-at-least-32-bytes"

// api รวมตัวทดสอบ connection pool ของ schema ชั่วคราว และ HTTP server ที่เรียกใช้งานจริง
type api struct {
	t      *testing.T
	pool   *pgxpool.Pool
	server *httptest.Server
}

// newAPI เตรียม schema/seed แยก เปิด HTTP server สำหรับทดสอบ และล้างทรัพยากรเมื่อจบ; ข้ามเมื่อไม่มี TEST_DATABASE_URL
func newAPI(t *testing.T) *api {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is required for PostgreSQL API integration tests")
	}
	ctx := context.Background()
	base, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal("parse test database configuration")
	}
	schema := fmt.Sprintf("api_test_%d", time.Now().UnixNano())
	if _, err = base.Exec(ctx, `CREATE SCHEMA `+schema); err != nil {
		base.Close()
		t.Fatal("create isolated test schema:", err)
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { pool.Close(); _, _ = base.Exec(ctx, `DROP SCHEMA `+schema+` CASCADE`); base.Close() })
	for i := 0; i < 2; i++ {
		for _, name := range []string{"schema.sql", "seed.sql"} {
			data, err := os.ReadFile(filepath.Join("..", name))
			if err != nil {
				t.Fatal(err)
			}
			if _, err = pool.Exec(ctx, string(data)); err != nil {
				t.Fatalf("%s repeat %d: %v", name, i, err)
			}
		}
	}
	h := &handler.Handler{Store: &store.Store{Pool: pool}, Tickets: &auth.Issuer{Secret: []byte(testSecret), TTL: time.Hour}}
	s := httptest.NewServer(mw.Recover(mw.CORS([]string{"http://localhost:3000", "https://thaifood.example"}, h.Routes())))
	t.Cleanup(s.Close)
	return &api{t: t, pool: pool, server: s}
}

// request ส่ง HTTP จริง ตรวจ status, JSON/error envelope และการไม่รั่ว password hash ก่อนคืนข้อมูลให้ test
func (a *api) request(method, path, token string, body any, want int) map[string]any {
	a.t.Helper()
	var data []byte
	if raw, ok := body.(string); ok {
		data = []byte(raw)
	} else if body != nil {
		data, _ = json.Marshal(body)
	}
	req, err := http.NewRequest(method, a.server.URL+path, bytes.NewReader(data))
	if err != nil {
		a.t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := a.server.Client().Do(req)
	if err != nil {
		a.t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		a.t.Fatal(err)
	}
	if resp.StatusCode != want {
		a.t.Fatalf("%s %s status %d, want %d: %s", method, path, resp.StatusCode, want, data)
	}
	if want == 204 {
		if len(data) != 0 {
			a.t.Fatalf("204 response has a body: %s", data)
		}
		return nil
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		a.t.Fatalf("non-JSON response to %s %s", method, path)
	}
	var out map[string]any
	if err = json.Unmarshal(data, &out); err != nil {
		a.t.Fatalf("invalid JSON: %v", err)
	}
	if want >= 400 {
		errObj, ok := out["error"].(map[string]any)
		if !ok || errObj["code"] == nil || errObj["message"] == nil {
			a.t.Fatalf("missing error envelope: %#v", out)
		}
	}
	if bytes.Contains(data, []byte("password_hash")) || bytes.Contains(data, []byte("passwordHash")) || bytes.Contains(data, []byte("$2a$")) || bytes.Contains(data, []byte("$2b$")) {
		a.t.Fatal("sensitive data leaked")
	}
	return out
}

// register สมัครสมาชิกและเข้าสู่ระบบเพื่อคืน ID กับ JWT สำหรับสถานการณ์ทดสอบถัดไป
func (a *api) register(email string) (int64, string) {
	a.t.Helper()
	m := a.request("POST", "/api/auth/register", "", map[string]any{"email": email, "password": "ExamplePass123!", "displayName": "สมาชิกทดสอบ"}, 201)
	if m["role"] != "member" {
		a.t.Fatalf("invalid register role: %v", m["role"])
	}
	l := a.request("POST", "/api/auth/login", "", map[string]any{"email": email, "password": "ExamplePass123!"}, 200)
	token, ok := l["token"].(string)
	if !ok || token == "" {
		a.t.Fatal("missing login JWT")
	}
	return int64(m["id"].(float64)), token
}

// TestAPIEndToEnd ตรวจข้อมูลสาธารณะ บัญชี ร้าน รีวิว ร้านโปรด สิทธิ์ คะแนนรวม และ cascade ผ่าน API จริง
func TestAPIEndToEnd(t *testing.T) {
	a := newAPI(t)
	a.request("GET", "/health", "", nil, 200)
	provinces := a.request("GET", "/api/provinces", "", nil, 200)["items"].([]any)
	if len(provinces) != 77 {
		t.Fatalf("got %d provinces, want77", len(provinces))
	}
	categories := a.request("GET", "/api/categories", "", nil, 200)["items"].([]any)
	if len(categories) == 0 {
		t.Fatal("missing categories")
	}
	list := a.request("GET", "/api/restaurants", "", nil, 200)
	if list["page"] != float64(1) || list["limit"] != float64(12) || list["total"] != float64(59) {
		t.Fatalf("pagination/seed mismatch: %#v", list)
	}
	for _, row := range list["items"].([]any) {
		r := row.(map[string]any)
		if r["rating"] != nil || r["reviewCount"] != float64(0) {
			t.Fatal("unreviewed seed restaurant has a rating")
		}
	}
	// ข้อมูลเริ่มต้นมีมากกว่าหนึ่งหน้า จึงต้องไล่หน้าตรวจว่าครบทุกร้านและไม่มี ID ซ้ำข้ามหน้า
	seen := make(map[float64]bool)
	covered := make(map[string]bool)
	for page, want := range []int{50, 9} {
		result := a.request("GET", fmt.Sprintf("/api/restaurants?limit=50&page=%d", page+1), "", nil, 200)
		rows := result["items"].([]any)
		if result["total"] != float64(59) || len(rows) != want {
			t.Fatalf("seed page %d has %d rows, want %d", page+1, len(rows), want)
		}
		for _, row := range rows {
			r := row.(map[string]any)
			id := r["id"].(float64)
			if seen[id] {
				t.Fatalf("restaurant %v appears on multiple pages", id)
			}
			seen[id] = true
			covered[r["province"].(string)] = true
		}
	}
	if len(seen) != 59 || len(covered) != 42 {
		t.Fatalf("seed pagination covers %d restaurants in %d provinces", len(seen), len(covered))
	}
	empty := a.request("GET", "/api/restaurants?q=__no_such_restaurant__", "", nil, 200)
	if len(empty["items"].([]any)) != 0 || empty["total"] != float64(0) {
		t.Fatal("empty search mismatch")
	}
	a.request("GET", "/api/restaurants?q=%27%20OR%201%3D1--", "", nil, 200)
	for _, query := range []string{"page=0", "page=-1", "page=bad", "limit=0", "limit=-1", "limit=bad"} {
		a.request("GET", "/api/restaurants?"+query, "", nil, 400)
	}
	if a.request("GET", "/api/restaurants?limit=100", "", nil, 200)["limit"] != float64(50) {
		t.Fatal("limit not capped50")
	}
	for _, path := range []string{"/api/restaurants/99999999", "/api/restaurants/99999999/reviews", "/does-not-exist"} {
		a.request("GET", path, "", nil, 404)
	}
	a.request("GET", "/api/restaurants/abc", "", nil, 400)
	a.request("POST", "/api/provinces", "", nil, 405)
	a.request("GET", "/api/me", "", nil, 401)
	a.request("POST", "/api/auth/register", "", `{"email":"hacker@example.com","password":"ExamplePass123!","displayName":"Hacker","role":"admin"}`, 400)
	for _, body := range []string{`{`, `[]`, `null`, `{} {}`, `{"email":"bad","password":"short","displayName":""}`} {
		a.request("POST", "/api/auth/register", "", body, 400)
	}
	// สร้างสมาชิกหลายบทบาทเพื่อตรวจการเข้าสู่ระบบ ข้อมูลส่วนตัว และสิทธิ์ที่แยกจากกัน
	id1, token1 := a.register("One@Example.com")
	_, token2 := a.register("two@example.com")
	adminID, adminToken := a.register("admin@example.com")
	a.request("POST", "/api/auth/register", "", map[string]any{"email": "one@example.com", "password": "ExamplePass123!", "displayName": "duplicate"}, 409)
	a.request("POST", "/api/auth/login", "", map[string]any{"email": "one@example.com", "password": "wrongpass"}, 401)
	a.request("POST", "/api/auth/login", "", map[string]any{"email": "missing@example.com", "password": "wrongpass"}, 401)
	me := a.request("GET", "/api/me", token1, nil, 200)
	if int64(me["id"].(float64)) != id1 || me["email"] != "one@example.com" {
		t.Fatal("wrong member identity")
	}
	updated := a.request("PATCH", "/api/me", token1, map[string]any{"displayName": "ชื่อใหม่"}, 200)
	if updated["displayName"] != "ชื่อใหม่" {
		t.Fatal("profile did not update")
	}
	a.request("PATCH", "/api/me", token1, map[string]any{"role": "admin"}, 400)
	a.request("PUT", "/api/me/password", token1, map[string]any{"currentPassword": "wrongpass", "newPassword": "ChangedPass123!"}, 401)
	a.request("PUT", "/api/me/password", token1, map[string]any{"currentPassword": "ExamplePass123!", "newPassword": "ChangedPass123!"}, 204)
	a.request("POST", "/api/auth/login", "", map[string]any{"email": "one@example.com", "password": "ExamplePass123!"}, 401)
	a.request("POST", "/api/auth/login", "", map[string]any{"email": "one@example.com", "password": "ChangedPass123!"}, 200)
	a.request("GET", "/api/admin/members", token1, nil, 403)
	if _, err := a.pool.Exec(context.Background(), `UPDATE members SET role='admin' WHERE id=$1`, adminID); err != nil {
		t.Fatal(err)
	}
	if len(a.request("GET", "/api/admin/members", adminToken, nil, 200)["items"].([]any)) != 3 {
		t.Fatal("admin member listing")
	}
	// ให้ admin สร้างร้านของ test เอง แล้วตรวจการแก้เฉพาะฟิลด์และการค้นหาชื่อไทย
	province := provinces[0].(map[string]any)
	category := categories[0].(map[string]any)
	input := map[string]any{"name": "ร้านทดสอบ API", "provinceId": province["id"], "categoryId": category["id"], "address": "ที่อยู่ทดสอบ", "description": "ร้านที่สร้างโดย integration test", "latitude": 13.75, "longitude": 100.5, "imageUrl": nil}
	a.request("POST", "/api/admin/restaurants", token1, input, 403)
	bad := map[string]any{"name": "bad", "provinceId": 99999999, "categoryId": category["id"], "address": "x", "description": "x"}
	a.request("POST", "/api/admin/restaurants", adminToken, bad, 400)
	created := a.request("POST", "/api/admin/restaurants", adminToken, input, 201)
	rid := int64(created["id"].(float64))
	path := fmt.Sprintf("/api/restaurants/%d", rid)
	if created["rating"] != nil || created["province"] != province["nameTh"] {
		t.Fatal("restaurant response mismatch")
	}
	a.request("GET", path, "", nil, 200)
	a.request("PUT", fmt.Sprintf("/api/admin/restaurants/%d", rid), adminToken, map[string]any{"name": "ร้านแก้ไข", "imageUrl": "https://example.com/image.jpg"}, 200)
	r := a.request("PUT", fmt.Sprintf("/api/admin/restaurants/%d", rid), adminToken, map[string]any{"imageUrl": nil}, 200)
	if r["name"] != "ร้านแก้ไข" || r["imageUrl"] != nil || r["address"] != input["address"] {
		t.Fatal("partial restaurant update lost fields")
	}
	filter := "/api/restaurants?" + "province=" + urlQuery(province["nameTh"].(string)) + "&category=" + urlQuery(category["nameTh"].(string)) + "&q=" + urlQuery("ร้านแก้ไข")
	if a.request("GET", filter, "", nil, 200)["total"] != float64(1) {
		t.Fatal("Thai name filters failed")
	}
	// ตรวจวงจรรีวิว ตั้งแต่ validation และ ownership จนถึงคะแนนเฉลี่ยที่เปลี่ยนตามรีวิวจริง
	if len(a.request("GET", path+"/reviews", "", nil, 200)["items"].([]any)) != 0 {
		t.Fatal("new restaurant reviews should be empty")
	}
	a.request("POST", path+"/reviews", "", map[string]any{"rating": 5, "text": "abcdefghij"}, 401)
	for _, body := range []map[string]any{{"rating": 0, "text": "abcdefghij"}, {"rating": 6, "text": "abcdefghij"}, {"rating": 1.5, "text": "abcdefghij"}, {"rating": 5, "text": "อร่อย"}} {
		a.request("POST", path+"/reviews", token1, body, 400)
	}
	a.request("POST", "/api/restaurants/99999999/reviews", token1, map[string]any{"rating": 5, "text": "abcdefghij"}, 404)
	review := a.request("POST", path+"/reviews", token1, map[string]any{"rating": 5, "text": "กขคงจฉชซฌญ"}, 201)
	review1 := int64(review["id"].(float64))
	if review["member"].(map[string]any)["displayName"] != "ชื่อใหม่" {
		t.Fatal("review author mismatch")
	}
	if _, err := time.Parse(time.RFC3339Nano, review["createdAt"].(string)); err != nil {
		t.Fatal("invalid review timestamp")
	}
	a.request("POST", path+"/reviews", token1, map[string]any{"rating": 4, "text": "abcdefghij"}, 409)
	review = a.request("POST", path+"/reviews", token2, map[string]any{"rating": 4, "text": "abcdefghij"}, 201)
	review2 := int64(review["id"].(float64))
	r = a.request("GET", path, "", nil, 200)
	if r["rating"] != 4.5 || r["reviewCount"] != float64(2) {
		t.Fatal("review aggregation failed")
	}
	a.request("DELETE", fmt.Sprintf("/api/reviews/%d", review1), token2, nil, 403)
	a.request("DELETE", fmt.Sprintf("/api/reviews/%d", review1), token1, nil, 204)
	r = a.request("GET", path, "", nil, 200)
	if r["rating"] != float64(4) || r["reviewCount"] != float64(1) {
		t.Fatal("rating not recomputed after owner deletion")
	}
	a.request("DELETE", fmt.Sprintf("/api/admin/reviews/%d", review2), token1, nil, 403)
	a.request("DELETE", fmt.Sprintf("/api/admin/reviews/%d", review2), adminToken, nil, 204)
	r = a.request("GET", path, "", nil, 200)
	if r["rating"] != nil || r["reviewCount"] != float64(0) {
		t.Fatal("no reviews must mean null rating")
	}
	a.request("DELETE", "/api/reviews/99999999", token1, nil, 404)
	// ตรวจว่าร้านโปรดแยกตามสมาชิก และการลบร้านล้างรีวิว/รายการโปรดที่อ้างถึงได้ครบ
	if len(a.request("GET", "/api/favorites", token1, nil, 200)["items"].([]any)) != 0 {
		t.Fatal("new member favorites should be empty")
	}
	a.request("POST", "/api/favorites", token1, map[string]any{"restaurantId": rid}, 201)
	a.request("POST", "/api/favorites", token1, map[string]any{"restaurantId": rid}, 409)
	a.request("POST", "/api/favorites", token1, map[string]any{"restaurantId": 99999999}, 404)
	if len(a.request("GET", "/api/favorites", token1, nil, 200)["items"].([]any)) != 1 {
		t.Fatal("favorite missing")
	}
	if len(a.request("GET", "/api/favorites", token2, nil, 200)["items"].([]any)) != 0 {
		t.Fatal("favorites leaked to another member")
	}
	a.request("DELETE", fmt.Sprintf("/api/favorites/%d", rid), token2, nil, 404)
	a.request("DELETE", fmt.Sprintf("/api/favorites/%d", rid), token1, nil, 204)
	a.request("POST", "/api/favorites", token1, map[string]any{"restaurantId": rid}, 201)
	a.request("POST", path+"/reviews", token1, map[string]any{"rating": 3, "text": "abcdefghij"}, 201)
	a.request("DELETE", fmt.Sprintf("/api/admin/restaurants/%d", rid), adminToken, nil, 204)
	a.request("GET", path, "", nil, 404)
	if len(a.request("GET", "/api/favorites", token1, nil, 200)["items"].([]any)) != 0 {
		t.Fatal("restaurant favorites not cascaded")
	}
	var count int
	if err := a.pool.QueryRow(context.Background(), `SELECT count(*) FROM reviews WHERE restaurant_id=$1`, rid).Scan(&count); err != nil || count != 0 {
		t.Fatal("reviews not cascaded")
	}
	// ใช้ token เดิมตรวจว่าการลดสิทธิ์มีผลทันที และ token ผิดรูปแบบ/หมดอายุยังถูกปฏิเสธ
	if _, err := a.pool.Exec(context.Background(), `UPDATE members SET role='member' WHERE id=$1`, adminID); err != nil {
		t.Fatal(err)
	}
	a.request("GET", "/api/admin/members", adminToken, nil, 403)
	for _, token := range []string{"invalid", expiredToken(t), wrongAlgorithmToken(t)} {
		a.request("GET", "/api/me", token, nil, 401)
	}
	t.Log("Verified all 20 Contract v2 endpoints + health, PostgreSQL persistence, role/ownership checks, validation, uniqueness and cascades")
}

// urlQuery แปลงแต่ละไบต์ของข้อความเป็น query encoding เพื่อส่งชื่อไทยและอักขระพิเศษได้ครบ
func urlQuery(s string) string { // keep UTF-8 names intact while encoding query separators.
	var b strings.Builder
	for _, c := range []byte(s) {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

// expiredToken สร้าง JWT ที่ลายเซ็นถูกต้องแต่หมดอายุ เพื่อทดสอบการปฏิเสธที่ middleware
func expiredToken(t *testing.T) string {
	t.Helper()
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "thaifood-guide", Audience: jwt.ClaimStrings{"thaifood-api"}, Subject: "1", IssuedAt: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour))}).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// wrongAlgorithmToken สร้าง JWT แบบ HS384 เพื่อยืนยันว่าระบบรับเฉพาะ HS256
func wrongAlgorithmToken(t *testing.T) string {
	t.Helper()
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.RegisteredClaims{Issuer: "thaifood-guide", Audience: jwt.ClaimStrings{"thaifood-api"}, Subject: "1", IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// TestConcurrentDuplicates ส่งคำขอเพิ่มรีวิว/ร้านโปรดซ้ำพร้อมกัน และตรวจว่ามี 201 หนึ่งครั้งกับ 409 ที่เหลือ
func TestConcurrentDuplicates(t *testing.T) {
	a := newAPI(t)
	_, token := a.register("concurrent@example.com")
	rows := a.request("GET", "/api/restaurants", "", nil, 200)["items"].([]any)
	rid := int64(rows[0].(map[string]any)["id"].(float64))
	for _, tc := range []struct {
		path string
		body any
	}{
		{fmt.Sprintf("/api/restaurants/%d/reviews", rid), map[string]any{"rating": 5, "text": "concurrent review"}},
		{"/api/favorites", map[string]any{"restaurantId": rid}},
	} {
		statuses := make(chan int, 6)
		var wg sync.WaitGroup
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, _ := json.Marshal(tc.body)
				req, _ := http.NewRequest("POST", a.server.URL+tc.path, bytes.NewReader(data))
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")
				resp, err := a.server.Client().Do(req)
				if err != nil {
					statuses <- 0
					return
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				statuses <- resp.StatusCode
			}()
		}
		wg.Wait()
		close(statuses)
		created, conflicts := 0, 0
		for status := range statuses {
			switch status {
			case 201:
				created++
			case 409:
				conflicts++
			default:
				t.Errorf("concurrent %s returned %d", tc.path, status)
			}
		}
		if created != 1 || conflicts != 5 {
			t.Fatalf("concurrent %s: created%d conflicts%d", tc.path, created, conflicts)
		}
	}
}
