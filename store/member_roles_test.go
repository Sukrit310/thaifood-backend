// ไฟล์ store/member_roles_test.go: ทดสอบ transaction เปลี่ยน role กับ PostgreSQL จริง รวมการแย่ง lock และการลดสิทธิ์พร้อมกัน
package store

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

// roleMember สร้างสมาชิกทดสอบและตั้ง role เริ่มต้นผ่าน SQL ใน schema ทดสอบเท่านั้น
func roleMember(t *testing.T, s *Store, email, role string) Member {
	t.Helper()
	m, err := s.CreateMember(context.Background(), email, "test-only-hash", "สมาชิกทดสอบ")
	if err != nil {
		t.Fatal(err)
	}
	if role == "admin" {
		if _, err := s.Pool.Exec(context.Background(), `UPDATE members SET role='admin' WHERE id=$1`, m.ID); err != nil {
			t.Fatal(err)
		}
		m.Role = role
	}
	return m
}

// TestUpdateMemberRolePersistenceAndIdempotence ตรวจการบันทึก role และการส่งค่าเดิมซ้ำโดยไม่เปลี่ยน updated_at
func TestUpdateMemberRolePersistenceAndIdempotence(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	actor := roleMember(t, s, "actor@example.com", "admin")
	target := roleMember(t, s, "target@example.com", "member")
	for _, role := range []string{"admin", "member"} {
		updated, err := s.UpdateMemberRole(ctx, actor.ID, target.ID, role)
		if err != nil {
			t.Fatal(err)
		}
		if updated.ID != target.ID || updated.Role != role || updated.Email != target.Email || updated.DisplayName != target.DisplayName || updated.PasswordHash != target.PasswordHash {
			t.Fatal("role update changed unrelated member fields or returned the wrong member")
		}
		stored, err := s.MemberByID(ctx, target.ID)
		if err != nil || stored != updated {
			t.Fatalf("role was not persisted: %v", err)
		}
		var before, after time.Time
		if err := s.Pool.QueryRow(ctx, `SELECT updated_at FROM members WHERE id=$1`, target.ID).Scan(&before); err != nil {
			t.Fatal(err)
		}
		again, err := s.UpdateMemberRole(ctx, actor.ID, target.ID, role)
		if err != nil || again != updated {
			t.Fatalf("same-role request was not idempotent: %v", err)
		}
		if err := s.Pool.QueryRow(ctx, `SELECT updated_at FROM members WHERE id=$1`, target.ID).Scan(&after); err != nil {
			t.Fatal(err)
		}
		if !before.Equal(after) {
			t.Fatal("same-role request changed updated_at")
		}
	}
}

// TestUpdateMemberRoleGuards ตรวจสิทธิ์ผู้กระทำ การแก้ตนเอง และการคืน lock หลัง transaction ถูกปฏิเสธ
func TestUpdateMemberRoleGuards(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	actor := roleMember(t, s, "actor@example.com", "admin")
	target := roleMember(t, s, "target@example.com", "member")
	for _, tt := range []struct {
		name              string
		actorID, targetID int64
		role              string
		want              error
	}{
		{"self demotion of last admin", actor.ID, actor.ID, "member", ErrForbidden},
		{"self unchanged", actor.ID, actor.ID, "admin", ErrForbidden},
		{"member caller", target.ID, actor.ID, "member", ErrForbidden},
		{"member no-op", target.ID, actor.ID, "admin", ErrForbidden},
		{"missing caller", 99999999, actor.ID, "member", ErrForbidden},
		{"missing target", actor.ID, 99999999, "admin", ErrNotFound},
		{"invalid role", actor.ID, target.ID, "superadmin", ErrForbidden},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.UpdateMemberRole(ctx, tt.actorID, tt.targetID, tt.role)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v, want %v", err, tt.want)
			}
			current, err := s.MemberByID(ctx, actor.ID)
			if err != nil || current.Role != "admin" {
				t.Fatalf("last admin changed after rejection: %v", err)
			}
			current, err = s.MemberByID(ctx, target.ID)
			if err != nil || current.Role != "member" {
				t.Fatalf("target changed after rejection: %v", err)
			}
		})
	}
	// transaction ที่ถูกปฏิเสธต้องคืน lock เพื่อไม่ขวางคำขอที่ถูกต้องในครั้งถัดไป
	if _, err := s.UpdateMemberRole(ctx, actor.ID, target.ID, "admin"); err != nil {
		t.Fatal("update after rejected requests:", err)
	}
}

// TestUpdateMemberRoleConcurrentCrossDemotions ให้ admin สองคนลดสิทธิ์กันพร้อมกัน แล้วตรวจว่าสำเร็จหนึ่งครั้งและยังเหลือ admin
func TestUpdateMemberRoleConcurrentCrossDemotions(t *testing.T) {
	s := testStore(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a := roleMember(t, s, "a@example.com", "admin")
	b := roleMember(t, s, "b@example.com", "admin")
	start := make(chan struct{})
	errs := make(chan error, 2)
	for _, pair := range [][2]int64{{a.ID, b.ID}, {b.ID, a.ID}} {
		go func() {
			<-start
			_, err := s.UpdateMemberRole(ctx, pair[0], pair[1], "member")
			errs <- err
		}()
	}
	close(start)
	successes, forbidden := 0, 0
	for range 2 {
		select {
		case err := <-errs:
			if err == nil {
				successes++
			} else if errors.Is(err, ErrForbidden) {
				forbidden++
			} else {
				t.Fatal(err)
			}
		case <-ctx.Done():
			t.Fatal("concurrent updates did not finish:", ctx.Err())
		}
	}
	if successes != 1 || forbidden != 1 {
		t.Fatalf("successes=%d forbidden=%d, want 1 each", successes, forbidden)
	}
	assertAdminCount(t, s, 1)
}

// TestUpdateMemberRoleRechecksActorAfterLockWait ตรวจว่าผู้ที่ถูกลดสิทธิ์ระหว่างรอ lock ไม่สามารถใช้สิทธิ์ admin เดิมได้
func TestUpdateMemberRoleRechecksActorAfterLockWait(t *testing.T) {
	s := testStore(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a := roleMember(t, s, "a@example.com", "admin")
	b := roleMember(t, s, "b@example.com", "admin")
	blocker, err := s.Pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer blocker.Rollback(context.Background())
	if _, err := blocker.Exec(ctx, `UPDATE members SET role='member' WHERE id=$1`, b.ID); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := s.UpdateMemberRole(ctx, b.ID, a.ID, "member")
		done <- err
	}()
	// รอจนฐานข้อมูลรายงานว่าติด lock จริง เพื่อทดสอบลำดับเหตุการณ์อย่างแน่นอน
	for {
		var waiting bool
		if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_locks WHERE relation='members'::regclass AND NOT granted)`).Scan(&waiting); err != nil {
			t.Fatal(err)
		}
		if waiting {
			break
		}
		select {
		case err := <-done:
			t.Fatalf("update did not wait for concurrent demotion: %v", err)
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(5 * time.Millisecond):
		}
	}
	if err := blocker.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("demoted caller retained admin authority: %v", err)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	assertAdminCount(t, s, 1)
}

// assertAdminCount นับ admin จากฐานข้อมูลจริงและเทียบกับจำนวนที่คาดหวัง
func assertAdminCount(t *testing.T, s *Store, want int) {
	t.Helper()
	var count int
	if err := s.Pool.QueryRow(context.Background(), `SELECT count(*) FROM members WHERE role='admin'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Fatal(fmt.Sprintf("admin count=%d, want %d", count, want))
	}
}
