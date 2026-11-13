// ไฟล์ integration/member_roles_test.go: ทดสอบสิทธิ์สมาชิกผ่าน HTTP และฐานข้อมูลจริง รวมผลต่อ JWT เดิมและคำขอที่เกิดพร้อมกัน
package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestAdminMemberRoles ตรวจเลื่อน/ลดสิทธิ์ การห้ามแก้ตนเอง validation และการที่ JWT เดิมใช้สิทธิ์ใหม่ได้ทันที
func TestAdminMemberRoles(t *testing.T) {
	a := newAPI(t)
	actorID, actorToken := a.register("actor@example.com")
	targetID, targetToken := a.register("target@example.com")
	if _, err := a.pool.Exec(context.Background(), `UPDATE members SET role='admin' WHERE id=$1`, actorID); err != nil {
		t.Fatal(err)
	}
	path := fmt.Sprintf("/api/admin/members/%d/role", targetID)
	a.request("PATCH", path, "", `{"role":"admin"}`, 401)
	a.request("PATCH", path, "invalid-token", `{"role":"admin"}`, 401)
	a.request("PATCH", path, targetToken, `{"role":"admin"}`, 403)
	a.request("GET", "/api/admin/members", targetToken, nil, 403)
	updated := a.request("PATCH", path, actorToken, `{"role":"admin"}`, 200)
	if len(updated) != 4 || updated["id"] != float64(targetID) || updated["email"] != "target@example.com" || updated["displayName"] != "สมาชิกทดสอบ" || updated["role"] != "admin" {
		t.Fatalf("unexpected Member response: %#v", updated)
	}
	if a.request("GET", "/api/me", targetToken, nil, 200)["role"] != "admin" {
		t.Fatal("promotion was not persisted")
	}
	// JWT เดิมต้องใช้สิทธิ์ admin ที่เพิ่งได้รับได้ทันที โดยไม่ต้องเข้าสู่ระบบใหม่
	a.request("GET", "/api/admin/members", targetToken, nil, 200)
	a.request("PATCH", path, actorToken, `{"role":"admin"}`, 200)
	for _, tokenAndID := range []struct {
		token string
		id    int64
	}{{actorToken, actorID}, {targetToken, targetID}} {
		for _, role := range []string{"member", "admin"} {
			a.request("PATCH", fmt.Sprintf("/api/admin/members/%d/role", tokenAndID.id), tokenAndID.token, map[string]string{"role": role}, 403)
		}
	}
	if a.request("PATCH", path, actorToken, `{"role":"member"}`, 200)["role"] != "member" {
		t.Fatal("demotion response did not contain the new role")
	}
	a.request("GET", "/api/admin/members", targetToken, nil, 403)
	a.request("PATCH", fmt.Sprintf("/api/admin/members/%d/role", actorID), targetToken, `{"role":"member"}`, 403)
	if a.request("GET", "/api/me", targetToken, nil, 200)["role"] != "member" {
		t.Fatal("demoted token cannot use ordinary member endpoints")
	}
	for _, role := range []string{"member", "admin"} {
		a.request("PATCH", fmt.Sprintf("/api/admin/members/%d/role", actorID), actorToken, map[string]string{"role": role}, 403)
	}
	for _, target := range []string{"0", "-1", "bad", "1.5", "9223372036854775808"} {
		a.request("PATCH", "/api/admin/members/"+target+"/role", actorToken, `{"role":"admin"}`, 400)
	}
	a.request("PATCH", "/api/admin/members/99999999/role", actorToken, `{"role":"admin"}`, 404)
	for _, payload := range []string{`{}`, `null`, `[]`, `{`, `{"role":null}`, `{"role":1}`, `{"role":"owner"}`, `{"role":"ADMIN"}`, `{"role":" admin "}`, `{"role":"admin","actorId":1}`, `{"Role":"admin"}`, `{"role":"admin"} {}`} {
		a.request("PATCH", path, actorToken, payload, 400)
	}
	a.request("PUT", path, actorToken, `{"role":"admin"}`, 405)
	a.request("POST", "/api/auth/register", "", `{"email":"escalate@example.com","password":"ExamplePass123!","displayName":"Test","role":"admin"}`, 400)
	a.request("PATCH", "/api/me", targetToken, `{"displayName":"New name","role":"admin"}`, 400)
	if a.request("GET", "/api/me", targetToken, nil, 200)["role"] != "member" {
		t.Fatal("invalid requests escalated role")
	}
	var count int
	if err := a.pool.QueryRow(context.Background(), `SELECT count(*) FROM members WHERE role='admin'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("last admin was not preserved: count=%d err=%v", count, err)
	}
}

// TestAdminMemberRolesConcurrentCrossDemotions ส่ง HTTP ให้ admin ลดสิทธิ์กันพร้อมกัน และตรวจว่ายังเหลือ admin หนึ่งคน
func TestAdminMemberRolesConcurrentCrossDemotions(t *testing.T) {
	a := newAPI(t)
	idA, tokenA := a.register("a@example.com")
	idB, tokenB := a.register("b@example.com")
	if _, err := a.pool.Exec(context.Background(), `UPDATE members SET role='admin' WHERE id IN ($1,$2)`, idA, idB); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	type result struct {
		status int
		err    error
	}
	results := make(chan result, 2)
	for _, attempt := range []struct {
		token  string
		target int64
	}{{tokenA, idB}, {tokenB, idA}} {
		go func() {
			<-start
			req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/admin/members/%d/role", a.server.URL, attempt.target), strings.NewReader(`{"role":"member"}`))
			if err != nil {
				results <- result{err: err}
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+attempt.token)
			resp, err := a.server.Client().Do(req)
			if err != nil {
				results <- result{err: err}
				return
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			results <- result{status: resp.StatusCode}
		}()
	}
	close(start)
	counts := map[int]int{}
	for range 2 {
		r := <-results
		if r.err != nil {
			t.Fatal(r.err)
		}
		counts[r.status]++
	}
	if counts[200] != 1 || counts[403] != 1 {
		t.Fatalf("concurrent responses: %v, want one 200 and one 403", counts)
	}
	var admins int
	if err := a.pool.QueryRow(context.Background(), `SELECT count(*) FROM members WHERE role='admin'`).Scan(&admins); err != nil || admins != 1 {
		t.Fatalf("concurrent demotions left %d admins: %v", admins, err)
	}
}
