package store

import (
	"context"
	"strings"
)

const memberColumns = `id, email, password_hash, display_name, role`

func (s *Store) CreateMember(ctx context.Context, email, hash, name string) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `INSERT INTO members(email,password_hash,display_name,role) VALUES($1,$2,$3,'member') RETURNING `+memberColumns, strings.ToLower(strings.TrimSpace(email)), hash, name))
}

func (s *Store) MemberByEmail(ctx context.Context, email string) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `SELECT `+memberColumns+` FROM members WHERE lower(email)=lower($1)`, strings.TrimSpace(email)))
}

func (s *Store) MemberByID(ctx context.Context, id int64) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `SELECT `+memberColumns+` FROM members WHERE id=$1`, id))
}

func (s *Store) UpdateMember(ctx context.Context, id int64, name string) (Member, error) {
	return scanMember(s.Pool.QueryRow(ctx, `UPDATE members SET display_name=$2,updated_at=now() WHERE id=$1 RETURNING `+memberColumns, id, name))
}

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
