package store

import "context"

func (s *Store) ListReviews(ctx context.Context, restaurantID int64) ([]Review, error) {
	var exists bool
	if err := s.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM restaurants WHERE id=$1)`, restaurantID).Scan(&exists); err != nil {
		return nil, dbError(err)
	}
	if !exists {
		return nil, ErrNotFound
	}
	rows, err := s.Pool.Query(ctx, `SELECT r.id,r.restaurant_id,r.rating,r.text,r.created_at,m.id,m.display_name FROM reviews r JOIN members m ON m.id=r.member_id WHERE r.restaurant_id=$1 ORDER BY r.created_at DESC,r.id DESC`, restaurantID)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	items := make([]Review, 0)
	for rows.Next() {
		r, err := scanReview(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, r)
	}
	return items, dbError(rows.Err())
}

func (s *Store) CreateReview(ctx context.Context, restaurantID, memberID int64, rating int, text string) (Review, error) {
	return scanReview(s.Pool.QueryRow(ctx, `WITH added AS (INSERT INTO reviews(restaurant_id,member_id,rating,text) VALUES($1,$2,$3,$4) RETURNING *)
	SELECT r.id,r.restaurant_id,r.rating,r.text,r.created_at,m.id,m.display_name FROM added r JOIN members m ON m.id=r.member_id`, restaurantID, memberID, rating, text))
}

func (s *Store) DeleteReview(ctx context.Context, id, actorID int64, admin bool) error {
	var exists, deleted bool
	err := s.Pool.QueryRow(ctx, `WITH existing AS (SELECT member_id FROM reviews WHERE id=$1),
	deleted AS (DELETE FROM reviews WHERE id=$1 AND (member_id=$2 OR $3) RETURNING id)
	SELECT EXISTS(SELECT 1 FROM existing),EXISTS(SELECT 1 FROM deleted)`, id, actorID, admin).Scan(&exists, &deleted)
	if err != nil {
		return dbError(err)
	}
	if deleted {
		return nil
	}
	if !exists {
		return ErrNotFound
	}
	return ErrForbidden
}
