package store

import "context"

func (s *Store) ListFavorites(ctx context.Context, memberID int64) ([]Restaurant, error) {
	rows, err := s.Pool.Query(ctx, restaurantSelect("restaurants")+` JOIN favorites f ON f.restaurant_id=r.id WHERE f.member_id=$1 ORDER BY f.created_at DESC,r.id DESC`, memberID)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	items := make([]Restaurant, 0)
	for rows.Next() {
		r, err := scanRestaurant(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, r)
	}
	return items, dbError(rows.Err())
}

func (s *Store) CreateFavorite(ctx context.Context, memberID, restaurantID int64) error {
	_, err := s.Pool.Exec(ctx, `INSERT INTO favorites(member_id,restaurant_id) VALUES($1,$2)`, memberID, restaurantID)
	return dbError(err)
}

func (s *Store) DeleteFavorite(ctx context.Context, memberID, restaurantID int64) error {
	tag, err := s.Pool.Exec(ctx, `DELETE FROM favorites WHERE member_id=$1 AND restaurant_id=$2`, memberID, restaurantID)
	return affected(tag, err)
}
