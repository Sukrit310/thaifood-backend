// ไฟล์ handler/member_roles_test.go: ทดสอบ endpoint เปลี่ยนสิทธิ์ด้วย repository จำลอง ครอบคลุม authorization, validation และ response
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"thaifood-backend/auth"
	"thaifood-backend/store"
)

// roleStore จำลองผู้กระทำ ผลลัพธ์ และการเรียกเปลี่ยน role โดยไม่ต้องเปิด PostgreSQL
type roleStore struct {
	store.Repository
	actor             store.Member
	updated           store.Member
	err               error
	calls             int
	actorID, targetID int64
	role              string
}

// MemberByID คืนผู้กระทำจำลองให้ middleware ตรวจสิทธิ์ก่อนเข้า endpoint
func (s *roleStore) MemberByID(context.Context, int64) (store.Member, error) {
	return s.actor, nil
}

// UpdateMemberRole บันทึกจำนวนครั้งและ arguments ที่ได้รับ แล้วคืนสมาชิกหรือ error ตามกรณีทดสอบ
func (s *roleStore) UpdateMemberRole(_ context.Context, actorID, targetID int64, role string) (store.Member, error) {
	s.calls++
	s.actorID, s.targetID, s.role = actorID, targetID, role
	return s.updated, s.err
}

// roleRequest สร้าง HTTP request สำหรับเปลี่ยน role พร้อมเลือก token, body และ Content-Type ที่ต้องการทดสอบ
func roleRequest(t *testing.T, s *roleStore, method, target, tokenKind, payload, media string) *httptest.ResponseRecorder {
	t.Helper()
	issuer := &auth.Issuer{Secret: []byte(strings.Repeat("x", 32))}
	r := httptest.NewRequest(method, "/api/admin/members/"+target+"/role", strings.NewReader(payload))
	r.Header.Set("Content-Type", media)
	if tokenKind != "guest" {
		token := "invalid-token"
		if tokenKind != "invalid" {
			var err error
			token, err = issuer.Issue(s.actor.ID)
			if err != nil {
				t.Fatal(err)
			}
		}
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	(&Handler{Store: s, Tickets: issuer}).Routes().ServeHTTP(w, r)
	return w
}

// TestUpdateMemberRoleAuthorization ตรวจว่า guest, token ผิด, member และผู้ที่แก้ role ตัวเองไม่ถึงขั้นเขียนฐานข้อมูล
func TestUpdateMemberRoleAuthorization(t *testing.T) {
	for _, tt := range []struct {
		name, token, actorRole, target string
		status                         int
		code                           string
	}{
		{"guest", "guest", "admin", "2", 401, "UNAUTHORIZED"},
		{"invalid token", "invalid", "admin", "2", 401, "UNAUTHORIZED"},
		{"member", "valid", "member", "2", 403, "FORBIDDEN"},
		{"self demotion", "valid", "admin", "1", 403, "FORBIDDEN"},
		{"self unchanged", "valid", "admin", "1", 403, "FORBIDDEN"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &roleStore{actor: store.Member{ID: 1, Role: tt.actorRole}}
			payload := `{"role":"member"}`
			if tt.name == "self unchanged" {
				payload = `{"role":"admin"}`
			}
			w := roleRequest(t, s, "PATCH", tt.target, tt.token, payload, "application/json")
			assertRoleError(t, w, tt.status, tt.code)
			if s.calls != 0 {
				t.Fatal("unauthorized request reached role update")
			}
		})
	}
}

// TestUpdateMemberRoleValidation ตรวจ ID, JSON, role และฟิลด์เกินก่อนเรียก repository
func TestUpdateMemberRoleValidation(t *testing.T) {
	for _, tt := range []struct{ name, target, payload, media string }{
		{"zero id", "0", `{"role":"admin"}`, "application/json"},
		{"negative id", "-1", `{"role":"admin"}`, "application/json"},
		{"text id", "abc", `{"role":"admin"}`, "application/json"},
		{"fractional id", "1.5", `{"role":"admin"}`, "application/json"},
		{"overflow id", "9223372036854775808", `{"role":"admin"}`, "application/json"},
		{"missing role", "2", `{}`, "application/json"},
		{"null role", "2", `{"role":null}`, "application/json"},
		{"empty role", "2", `{"role":""}`, "application/json"},
		{"unknown role", "2", `{"role":"owner"}`, "application/json"},
		{"uppercase role", "2", `{"role":"ADMIN"}`, "application/json"},
		{"padded role", "2", `{"role":" admin "}`, "application/json"},
		{"numeric role", "2", `{"role":1}`, "application/json"},
		{"boolean role", "2", `{"role":true}`, "application/json"},
		{"object role", "2", `{"role":{}}`, "application/json"},
		{"array role", "2", `{"role":["admin"]}`, "application/json"},
		{"unknown field", "2", `{"role":"admin","actorId":1}`, "application/json"},
		{"cased field", "2", `{"Role":"admin"}`, "application/json"},
		{"extra cased field", "2", `{"role":"admin","ROLE":"member"}`, "application/json"},
		{"malformed", "2", `{"role":`, "application/json"},
		{"empty body", "2", ``, "application/json"},
		{"null body", "2", `null`, "application/json"},
		{"array body", "2", `[]`, "application/json"},
		{"trailing JSON", "2", `{"role":"admin"} {}`, "application/json"},
		{"trailing garbage", "2", `{"role":"admin"} x`, "application/json"},
		{"wrong content type", "2", `{"role":"admin"}`, "text/plain"},
		{"oversized body", "2", `{"role":"` + strings.Repeat("x", 1<<20) + `"}`, "application/json"},
		{"invalid UTF8", "2", "{\"role\":\"\xff\"}", "application/json"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &roleStore{actor: store.Member{ID: 1, Role: "admin"}}
			w := roleRequest(t, s, "PATCH", tt.target, "valid", tt.payload, tt.media)
			assertRoleError(t, w, 400, "VALIDATION_ERROR")
			if s.calls != 0 {
				t.Fatal("invalid request reached role update")
			}
		})
	}
}

// TestUpdateMemberRoleReturnsMember ตรวจ arguments ของการเปลี่ยนสิทธิ์และ response สมาชิกที่ไม่มี password hash
func TestUpdateMemberRoleReturnsMember(t *testing.T) {
	for _, role := range []string{"admin", "member"} {
		t.Run(role, func(t *testing.T) {
			s := &roleStore{
				actor:   store.Member{ID: 7, Role: "admin"},
				updated: store.Member{ID: 2, Email: "target@example.com", DisplayName: "ชื่อสมาชิก", Role: role, PasswordHash: "private-hash"},
			}
			w := roleRequest(t, s, "PATCH", "2", "valid", `{"role":"`+role+`"}`, "application/json; charset=utf-8")
			if w.Code != 200 {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if s.calls != 1 || s.actorID != 7 || s.targetID != 2 || s.role != role {
				t.Fatalf("wrong authenticated identity or update arguments: %+v", s)
			}
			var out map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
				t.Fatal(err)
			}
			if len(out) != 4 || out["id"] != float64(2) || out["email"] != "target@example.com" || out["displayName"] != "ชื่อสมาชิก" || out["role"] != role {
				t.Fatalf("wrong Member response: %s", w.Body.String())
			}
			if !strings.Contains(w.Header().Get("Content-Type"), "application/json") || strings.Contains(w.Body.String(), "private-hash") {
				t.Fatalf("unsafe response: %s", w.Body.String())
			}
		})
	}
}

// TestUpdateMemberRoleStoreErrors ตรวจสถานะและรหัส error เมื่อ repository ปฏิเสธการเปลี่ยนสิทธิ์
func TestUpdateMemberRoleStoreErrors(t *testing.T) {
	for _, tt := range []struct {
		err    error
		status int
		code   string
	}{
		{store.ErrNotFound, 404, "NOT_FOUND"},
		{store.ErrForbidden, 403, "FORBIDDEN"},
		{store.ErrConflict, 409, "CONFLICT"},
		{store.ErrLastAdmin, 409, "LAST_ADMIN"},
		{errors.New("database password=private-secret"), 500, "INTERNAL_ERROR"},
	} {
		t.Run(tt.code, func(t *testing.T) {
			s := &roleStore{actor: store.Member{ID: 1, Role: "admin"}, err: tt.err}
			w := roleRequest(t, s, "PATCH", "2", "valid", `{"role":"admin"}`, "application/json")
			assertRoleError(t, w, tt.status, tt.code)
			if strings.Contains(w.Body.String(), "private-secret") {
				t.Fatal("database details leaked")
			}
		})
	}
}

// assertRoleError ตรวจ HTTP status และโครงสร้าง JSON error พร้อมรหัสที่คาดหวัง
func assertRoleError(t *testing.T, w *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	var out struct {
		Error struct{ Code, Message string }
	}
	if w.Code != status || json.Unmarshal(w.Body.Bytes(), &out) != nil || out.Error.Code != code || out.Error.Message == "" {
		t.Fatalf("want %d %s, got %d %s", status, code, w.Code, w.Body.String())
	}
}
