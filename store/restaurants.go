// ไฟล์ store/restaurants.go: อ่าน/เขียนร้านและข้อมูลอ้างอิง พร้อมค้นหา แบ่งหน้า และคำนวณคะแนนจากรีวิว
package store

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

// ListProvinces อ่านจังหวัดเรียงตามชื่อไทยและ ID เพื่อให้ผลลัพธ์มีลำดับแน่นอน
func (s *Store) ListProvinces(ctx context.Context) ([]Province, error) {
	rows, err := s.Pool.Query(ctx, `SELECT id,name_th,region FROM provinces ORDER BY name_th,id`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	items := make([]Province, 0)
	for rows.Next() {
		var p Province
		if err := rows.Scan(&p.ID, &p.NameTh, &p.Region); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, dbError(rows.Err())
}

// ListCategories อ่านประเภทอาหารเรียงตาม ID
func (s *Store) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := s.Pool.Query(ctx, `SELECT id,name_th FROM categories ORDER BY id`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	items := make([]Category, 0)
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.NameTh); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, dbError(rows.Err())
}

// from ต้องเป็นชื่อตาราง/CTE ที่โค้ดกำหนดเองเท่านั้น เพราะถูกนำมาต่อเป็น SQL identifier โดยตรง
// restaurantSelect สร้าง SELECT ที่ join จังหวัด/ประเภทอาหารและคำนวณค่าเฉลี่ยกับจำนวนรีวิวของแต่ละร้าน
func restaurantSelect(from string) string {
	return `SELECT r.id,r.name,p.name_th,c.name_th,r.address,r.description,r.latitude,r.longitude,r.image_url,stats.rating,stats.review_count
	FROM ` + from + ` r JOIN provinces p ON p.id=r.province_id JOIN categories c ON c.id=r.category_id
	LEFT JOIN LATERAL (SELECT avg(v.rating)::double precision AS rating,count(*) AS review_count FROM reviews v WHERE v.restaurant_id=r.id) stats ON true`
}

// เงื่อนไขร่วมของคำสั่งนับและอ่านหน้า: สตริงว่างหมายถึงไม่กรอง และใช้ parameters แทนการต่อค่าจากผู้ใช้
const restaurantWhere = ` WHERE ($1='' OR p.name_th=$1) AND ($2='' OR c.name_th=$2)
	AND ($3='' OR r.name ILIKE $3 ESCAPE '\' OR r.description ILIKE $3 ESCAPE '\' OR r.address ILIKE $3 ESCAPE '\')`

// searchPattern escape เครื่องหมาย wildcard ของ LIKE เพื่อให้ค้นข้อความตามตัวอักษรจริง
func searchPattern(query string) string {
	if query == "" {
		return ""
	}
	return "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(query) + "%"
}

// ListRestaurants ใช้ transaction แบบ repeatable read ให้จำนวนรวมและข้อมูลหน้าที่อ่านมาจาก snapshot เดียวกัน
func (s *Store) ListRestaurants(ctx context.Context, f RestaurantFilter) ([]Restaurant, int64, error) {
	// อ่าน count และรายการภายใต้ snapshot เดียว เพื่อให้ total ไม่คลาดกับหน้าที่อ่านเมื่อมีการเขียนพร้อมกัน
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, 0, dbError(err)
	}
	defer tx.Rollback(ctx)
	args := []any{f.Province, f.Category, searchPattern(f.Query)}
	var total int64
	err = tx.QueryRow(ctx, `SELECT count(*) FROM restaurants r JOIN provinces p ON p.id=r.province_id JOIN categories c ON c.id=r.category_id`+restaurantWhere, args...).Scan(&total)
	if err != nil {
		return nil, 0, dbError(err)
	}
	args = append(args, f.Limit, int64(f.Page-1)*int64(f.Limit))
	rows, err := tx.Query(ctx, restaurantSelect("restaurants")+restaurantWhere+` ORDER BY r.id LIMIT $4 OFFSET $5`, args...)
	if err != nil {
		return nil, 0, dbError(err)
	}
	items := make([]Restaurant, 0)
	for rows.Next() {
		r, err := scanRestaurant(rows)
		if err != nil {
			rows.Close()
			return nil, 0, err
		}
		items = append(items, r)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, 0, dbError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, 0, dbError(err)
	}
	return items, total, nil
}

// RestaurantByID อ่านรายละเอียดร้านหนึ่งรายการพร้อมสถิติรีวิว
func (s *Store) RestaurantByID(ctx context.Context, id int64) (Restaurant, error) {
	return scanRestaurant(s.Pool.QueryRow(ctx, restaurantSelect("restaurants")+` WHERE r.id=$1`, id))
}

// CreateRestaurant เพิ่มร้านด้วย SQL parameters และอ่านรายละเอียดที่ join แล้วจากแถวที่เพิ่งเพิ่ม
func (s *Store) CreateRestaurant(ctx context.Context, in RestaurantInput) (Restaurant, error) {
	return scanRestaurant(s.Pool.QueryRow(ctx, `WITH changed AS (
	INSERT INTO restaurants(name,province_id,category_id,address,description,latitude,longitude,image_url)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *) `+restaurantSelect("changed"), in.Name, in.ProvinceID, in.CategoryID, in.Address, in.Description, in.Latitude, in.Longitude, in.ImageURL))
}

// UpdateRestaurant ใช้ coalesce คงค่าฟิลด์ที่ไม่ได้ส่ง และใช้ ImageURLSet แยกการล้างรูปออกจากการไม่แก้รูป
func (s *Store) UpdateRestaurant(ctx context.Context, id int64, in RestaurantInput) (Restaurant, error) {
	return scanRestaurant(s.Pool.QueryRow(ctx, `WITH changed AS (
	UPDATE restaurants SET name=coalesce($2,name),province_id=coalesce($3,province_id),category_id=coalesce($4,category_id),
	address=coalesce($5,address),description=coalesce($6,description),latitude=coalesce($7,latitude),longitude=coalesce($8,longitude),
	image_url=CASE WHEN $9 THEN $10::text ELSE image_url END,updated_at=now() WHERE id=$1 RETURNING *) `+restaurantSelect("changed"),
		id, in.Name, in.ProvinceID, in.CategoryID, in.Address, in.Description, in.Latitude, in.Longitude, in.ImageURLSet, in.ImageURL))
}

// DeleteRestaurant ลบร้านตาม ID และแจ้ง ErrNotFound หากไม่พบร้านนั้น
func (s *Store) DeleteRestaurant(ctx context.Context, id int64) error {
	tag, err := s.Pool.Exec(ctx, `DELETE FROM restaurants WHERE id=$1`, id)
	return affected(tag, err)
}
