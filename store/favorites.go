// ไฟล์ store/favorites.go: อ่าน เพิ่ม และลบร้านโปรดโดยจำกัดข้อมูลตามรหัสสมาชิก
package store

import "context"

// ListFavorites อ่านร้านโปรดพร้อมรายละเอียดและคะแนน เรียงรายการที่เพิ่มล่าสุดก่อน
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

// CreateFavorite เพิ่มคู่สมาชิก/ร้าน โดยให้ primary key ป้องกันการเพิ่มร้านเดิมซ้ำ
func (s *Store) CreateFavorite(ctx context.Context, memberID, restaurantID int64) error {
	_, err := s.Pool.Exec(ctx, `INSERT INTO favorites(member_id,restaurant_id) VALUES($1,$2)`, memberID, restaurantID)
	return dbError(err)
}

// DeleteFavorite ลบเฉพาะคู่สมาชิก/ร้านที่ระบุ และแจ้ง ErrNotFound เมื่อยังไม่ได้บันทึกร้านนี้
func (s *Store) DeleteFavorite(ctx context.Context, memberID, restaurantID int64) error {
	tag, err := s.Pool.Exec(ctx, `DELETE FROM favorites WHERE member_id=$1 AND restaurant_id=$2`, memberID, restaurantID)
	return affected(tag, err)
}
