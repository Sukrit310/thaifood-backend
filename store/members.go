// ไฟล์ store/members.go: คำสั่งฐานข้อมูลสำหรับสมาชิก โปรไฟล์ การเปลี่ยนรหัสผ่าน และการเปลี่ยนสิทธิ์แบบ transaction
package store

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

// ลำดับคอลัมน์นี้ต้องตรงกับ scanMember ทุกครั้งที่อ่านหรือคืนข้อมูลสมาชิก
const memberColumns = `id, email, password_hash, display_name, role`

// CreateMember ปรับอีเมลเป็นตัวเล็กและสร้างบัญชีด้วย role member เสมอ
func (s *Store) CreateMember(ctx context.Context, email, hash, name string) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `INSERT INTO members(email,password_hash,display_name,role) VALUES($1,$2,$3,'member') RETURNING `+memberColumns, strings.ToLower(strings.TrimSpace(email)), hash, name))
}

// MemberByEmail ค้นสมาชิกโดยเปรียบเทียบอีเมลแบบไม่แยกตัวพิมพ์เล็กใหญ่
func (s *Store) MemberByEmail(ctx context.Context, email string) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `SELECT `+memberColumns+` FROM members WHERE lower(email)=lower($1)`, strings.TrimSpace(email)))
}

// MemberByID โหลดสมาชิกพร้อม role และ password hash ล่าสุดตามรหัสสมาชิก
func (s *Store) MemberByID(ctx context.Context, id int64) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `SELECT `+memberColumns+` FROM members WHERE id=$1`, id))
}

// UpdateMember เปลี่ยนชื่อแสดงผลและเวลาที่แก้ไข แล้วคืนข้อมูลสมาชิกหลังบันทึก
func (s *Store) UpdateMember(ctx context.Context, id int64, name string) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `UPDATE members SET display_name=$2,updated_at=now() WHERE id=$1 RETURNING `+memberColumns, id, name))
}

// UpdateMemberRole ล็อกการเขียนตารางสมาชิก ตรวจสิทธิ์ผู้กระทำซ้ำ และป้องกันการลดสิทธิ์ admin คนสุดท้าย
func (s *Store) UpdateMemberRole(ctx context.Context, actorID, targetID int64, role string) (Member, error) {
	if actorID == targetID || (role != "member" && role != "admin") {
		return Member{}, ErrForbidden
	}
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return Member{}, dbError(err)
	}
	defer tx.Rollback(ctx)
	// ล็อกการเขียนตารางสมาชิกข้าม API instances โดยยังอ่านปกติได้
	// READ COMMITTED ทำให้การตรวจหลังรอ lock เห็นข้อมูลที่ transaction อื่นเพิ่ง commit
	if _, err := tx.Exec(ctx, `LOCK TABLE members IN SHARE ROW EXCLUSIVE MODE`); err != nil {
		return Member{}, dbError(err)
	}
	// ตรวจสิทธิ์ผู้กระทำอีกครั้งหลังได้ lock เพราะระหว่างรออาจถูก admin คนอื่นลดสิทธิ์แล้ว
	var actorRole string
	err = tx.QueryRow(ctx, `SELECT role FROM members WHERE id=$1`, actorID).Scan(&actorRole)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && actorRole != "admin") {
		return Member{}, ErrForbidden
	}
	if err != nil {
		return Member{}, dbError(err)
	}
	target, err := scanMember(tx.QueryRow(ctx, `SELECT `+memberColumns+` FROM members WHERE id=$1`, targetID))
	if err != nil {
		return Member{}, err
	}
	if target.Role != role {
		if target.Role == "admin" && role == "member" {
			var admins int64
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM members WHERE role='admin'`).Scan(&admins); err != nil {
				return Member{}, dbError(err)
			}
			// ตรวจจำนวน admin เพิ่มอีกชั้น นอกเหนือจากการห้ามแก้ตนเองและตรวจสิทธิ์ผู้กระทำ
			if admins <= 1 {
				return Member{}, ErrLastAdmin
			}
		}
		target, err = scanMember(tx.QueryRow(ctx, `UPDATE members SET role=$2,updated_at=now() WHERE id=$1 RETURNING `+memberColumns, targetID, role))
		if err != nil {
			return Member{}, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return Member{}, dbError(err)
	}
	return target, nil
}

// UpdatePassword บันทึก hash ใหม่เมื่อ hash เดิมยังตรงกัน เพื่อไม่เขียนทับการเปลี่ยนรหัสผ่านที่เกิดพร้อมกัน
func (s *Store) UpdatePassword(ctx context.Context, id int64, oldHash, newHash string) error {
	tag, err := s.Pool.Exec(ctx, `UPDATE members SET password_hash=$3,updated_at=now() WHERE id=$1 AND password_hash=$2`, id, oldHash, newHash)
	if err != nil {
		return dbError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}

// ListMembers อ่านสมาชิกทั้งหมดเรียงตาม ID โดย PasswordHash จะถูกซ่อนเมื่อแปลงเป็น JSON
func (s *Store) ListMembers(ctx context.Context) ([]Member, error) {
	rows, err := s.Pool.Query(ctx, `SELECT `+memberColumns+` FROM members ORDER BY id`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	items := make([]Member, 0)
	for rows.Next() {
		m, err := scanMember(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, dbError(rows.Err())
}
