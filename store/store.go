package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct{ Pool *pgxpool.Pool }

var _ Repository = (*Store)(nil)

func (s *Store) Ping(ctx context.Context) error { return s.Pool.Ping(ctx) }

func dbError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var e *pgconn.PgError
	if errors.As(err, &e) {
		switch e.Code {
		case "23505":
			return ErrConflict
		case "23503", "23502", "23514":
			return ErrForeignKey
		}
	}
	return err
}

type scanner interface{ Scan(...any) error }

func scanMember(row scanner) (Member, error) {
	var m Member
	err := row.Scan(&m.ID, &m.Email, &m.PasswordHash, &m.DisplayName, &m.Role)
	return m, dbError(err)
}

func scanRestaurant(row scanner) (Restaurant, error) {
	var r Restaurant
	err := row.Scan(&r.ID, &r.Name, &r.Province, &r.Category, &r.Address, &r.Description, &r.Latitude, &r.Longitude, &r.ImageURL, &r.Rating, &r.ReviewCount)
	return r, dbError(err)
}

func scanReview(row scanner) (Review, error) {
	var r Review
	err := row.Scan(&r.ID, &r.RestaurantID, &r.Rating, &r.Text, &r.CreatedAt, &r.Member.ID, &r.Member.DisplayName)
	r.CreatedAt = r.CreatedAt.UTC()
	return r, dbError(err)
}

func affected(tag pgconn.CommandTag, err error) error {
	if err != nil {
		return dbError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
