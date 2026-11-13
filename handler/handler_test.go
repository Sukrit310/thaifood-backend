// ไฟล์ handler/handler_test.go: ใช้ repository จำลองทดสอบ routing, pagination, validation และสิทธิ์ โดยไม่เชื่อมฐานข้อมูล
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"thaifood-backend/auth"
	"thaifood-backend/store"
)

// testStore เก็บข้อมูลที่ส่งเข้า repository และกำหนดผลตอบกลับสำหรับ unit test ของ handler
type testStore struct {
	store.Repository
	filter                   store.RestaurantFilter
	member                   store.Member
	err                      error
	email, hash, displayName string
	rating                   int
	reviewText               string
}

// Ping จำลองผลการตรวจฐานข้อมูลตาม error ที่แต่ละกรณีทดสอบกำหนด
func (s *testStore) Ping(context.Context) error { return s.err }

// ListRestaurants เก็บ filter ที่ handler ส่งมาและคืนรายการว่างสำหรับตรวจ pagination
func (s *testStore) ListRestaurants(_ context.Context, f store.RestaurantFilter) ([]store.Restaurant, int64, error) {
	s.filter = f
	return nil, 0, s.err
}

// MemberByID คืนสมาชิกจำลองให้ middleware ใช้ตรวจตัวตนและสิทธิ์
func (s *testStore) MemberByID(context.Context, int64) (store.Member, error) { return s.member, s.err }

// MemberByEmail คืนสมาชิกหรือ error จำลองสำหรับทดสอบการเข้าสู่ระบบ
func (s *testStore) MemberByEmail(context.Context, string) (store.Member, error) {
	return s.member, s.err
}

// CreateMember เก็บข้อมูลสมัครที่ได้รับ เพื่อตรวจการปรับอีเมล ชื่อ และ password hash
func (s *testStore) CreateMember(_ context.Context, email, hash, displayName string) (store.Member, error) {
	s.email, s.hash, s.displayName = email, hash, displayName
	return store.Member{ID: 1, Email: email, PasswordHash: hash, DisplayName: displayName, Role: "member"}, s.err
}

// CreateReview เก็บคะแนนและข้อความที่ผ่าน validation แล้ว เพื่อตรวจค่าที่จะถูกบันทึก
func (s *testStore) CreateReview(_ context.Context, restaurantID, memberID int64, rating int, text string) (store.Review, error) {
	s.rating, s.reviewText = rating, text
	return store.Review{ID: 1, RestaurantID: restaurantID, Rating: rating, Text: text, Member: store.ReviewMember{ID: memberID, DisplayName: s.member.DisplayName}}, s.err
}

// TestPublicErrorsAndPagination ตรวจสถานะและรูปแบบ error ค่า pagination และการไม่เปิดเผยรายละเอียดฐานข้อมูล
func TestPublicErrorsAndPagination(t *testing.T) {
	s := &testStore{}
	h := (&Handler{Store: s, Tickets: &auth.Issuer{Secret: []byte(strings.Repeat("x", 32))}}).Routes()
	for _, tt := range []struct {
		method, path, body string
		status             int
	}{
		{"GET", "/unknown", "", 404}, {"POST", "/api/restaurants", "", 405},
		{"GET", "/api/restaurants?page=0", "", 400}, {"GET", "/api/restaurants?limit=no", "", 400},
		{"GET", "/api/restaurants/-1", "", 400},
		{"POST", "/api/auth/register", `{"email":"test@example.com","password":"12345678","displayName":"Test","role":"admin"}`, 400},
		{"POST", "/api/auth/register", `null`, 400},
		{"POST", "/api/auth/register", `{"email":"test@example.com","password":"12345678","displayName":"Test"} {}`, 400},
		{"POST", "/api/auth/register", `{"email":"test@example.com","password":"123","displayName":"Test"}`, 400},
		{"POST", "/api/auth/register", `{"email":"test@example.com","password":"12345678","displayName":"Test\u0000Name"}`, 400},
		{"GET", "/api/me", "", 401}, {"GET", "/api/admin/members", "", 401},
	} {
		t.Run(tt.method+tt.path+tt.body, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
			if w.Code != tt.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			var result struct {
				Error struct{ Code, Message string }
			}
			if json.Unmarshal(w.Body.Bytes(), &result) != nil || result.Error.Code == "" || result.Error.Message == "" {
				t.Fatalf("invalid envelope %s", w.Body.String())
			}
		})
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/restaurants?limit=100&province=%E0%B9%84%E0%B8%97%E0%B8%A2", nil))
	if w.Code != 200 || s.filter.Limit != 50 || s.filter.Page != 1 || !strings.Contains(w.Body.String(), `"items":[]`) {
		t.Fatalf("bad pagination: %+v %s", s.filter, w.Body.String())
	}
	s.err = errors.New("postgresql://secret:password@host/db")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/restaurants", nil))
	if w.Code != 500 || strings.Contains(w.Body.String(), "password") {
		t.Fatalf("leaked database error: %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/health", nil))
	if w.Code != 503 {
		t.Fatalf("health=%d", w.Code)
	}
}

// TestRestaurantInputValidation ตรวจฟิลด์ร้านที่ผิดรูปแบบและความหมายของ imageUrl ที่ส่งเป็น null
func TestRestaurantInputValidation(t *testing.T) {
	for _, body := range []string{`{}`, `{"name":null}`, `{"latitude":91,"longitude":100}`, `{"latitude":10}`, `{"imageUrl":"javascript:alert(1)"}`, `{"provinceId":0}`, `{"role":"admin"}`} {
		r := httptest.NewRequest("PUT", "/", strings.NewReader(body))
		w := httptest.NewRecorder()
		_, err := restaurantInput(w, r, false)
		if err == nil {
			t.Errorf("accepted %s", body)
		}
	}
	r := httptest.NewRequest("PUT", "/", strings.NewReader(`{"imageUrl":null}`))
	in, err := restaurantInput(httptest.NewRecorder(), r, false)
	if err != nil || !in.ImageURLSet || in.ImageURL != nil {
		t.Fatalf("explicit null: %+v %v", in, err)
	}
}

// TestMemberCannotUseAdminRoute ตรวจว่า JWT ของ member เข้า endpoint สำหรับ admin แล้วได้ 403
func TestMemberCannotUseAdminRoute(t *testing.T) {
	issuer := &auth.Issuer{Secret: []byte(strings.Repeat("x", 32))}
	token, err := issuer.Issue(1)
	if err != nil {
		t.Fatal(err)
	}
	s := &testStore{member: store.Member{ID: 1, Role: "member"}}
	h := (&Handler{Store: s, Tickets: issuer}).Routes()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/admin/members", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}
