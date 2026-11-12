package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is required for PostgreSQL integration tests")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("thaifood_store_%d", time.Now().UnixNano())
	if _, err = admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE"); admin.Close() })
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	for _, name := range []string{"schema.sql", "seed.sql", "schema.sql", "seed.sql"} {
		b, err := os.ReadFile(filepath.Join("..", name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, string(b)); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	return &Store{Pool: pool}
}

func TestSeedIsIdempotentAndHasAllProvinces(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	p, err := s.ListProvinces(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 77 {
		t.Fatalf("got %d provinces, want 77", len(p))
	}
	items, total, err := s.ListRestaurants(ctx, RestaurantFilter{Page: 1, Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	if total != 30 || len(items) != 30 {
		t.Fatalf("seed duplicate/missing restaurants: %d/%d", total, len(items))
	}
	covered := make(map[string]bool)
	for _, r := range items {
		covered[r.Province] = true
		if r.Rating != nil || r.ReviewCount != 0 {
			t.Fatalf("unreviewed restaurant has rating: %+v", r)
		}
	}
	if len(covered) != 19 {
		t.Fatalf("seed covers %d provinces, want 19 after adding ten previously empty provinces", len(covered))
	}
	var unsecured int
	if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM pg_class WHERE relnamespace=current_schema()::regnamespace AND relkind='r' AND NOT relrowsecurity`).Scan(&unsecured); err != nil {
		t.Fatal(err)
	}
	if unsecured != 0 {
		t.Fatalf("%d tables missing RLS", unsecured)
	}
}

func TestSeedProvidesRealRestaurantsWithAttributedPhotos(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	items, _, err := s.ListRestaurants(ctx, RestaurantFilter{Page: 1, Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range items {
		if r.ImageURL == nil {
			t.Fatalf("restaurant %q has no real photograph", r.Name)
		}
		u, err := url.Parse(*r.ImageURL)
		if err != nil || u.Scheme != "https" || (u.Host != "thumb.wikimedia.org" && u.Host != "upload.wikimedia.org") {
			t.Fatalf("restaurant %q has an invalid photograph URL", r.Name)
		}
		if strings.Contains(r.Name, "สมมติ") || strings.Contains(r.Address, "สมมติ") {
			t.Fatalf("fictional restaurant remains in a fresh seed: %q", r.Name)
		}
		for _, required := range []string{"ที่มาข้อมูล:", "ภาพ:", "https://commons.wikimedia.org/", "https://creativecommons.org/"} {
			if !strings.Contains(r.Description, required) {
				t.Fatalf("restaurant %q is missing data/image attribution %q", r.Name, required)
			}
		}
	}
}

func TestReseedRetiresOnlyUntouchedUnreferencedFictionalRestaurants(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	m, err := s.CreateMember(ctx, "seed-preservation@example.com", "test-only-hash", "Seed preservation")
	if err != nil {
		t.Fatal(err)
	}
	ids := make([]int64, 4)
	for i := range ids {
		err := s.Pool.QueryRow(ctx, `INSERT INTO restaurants(name,province_id,category_id,address,description,source_name,source_url)
			SELECT 'ร้านตัวอย่างสมมติ ทดสอบ',p.id,c.id,'ที่อยู่สมมติ','ข้อมูลสมมติ',
			'ThaiFood Guide author-created fictional sample',$1
			FROM provinces p CROSS JOIN categories c WHERE p.name_th='เชียงใหม่' AND c.name_th='อาหารเหนือ' RETURNING id`,
			fmt.Sprintf("urn:thaifood-guide:fictional-sample:test:%d", i)).Scan(&ids[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.CreateReview(ctx, ids[1], m.ID, 5, "รีวิวนี้ต้องไม่ถูกลบเมื่อปรับชุดข้อมูล"); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateFavorite(ctx, m.ID, ids[2]); err != nil {
		t.Fatal(err)
	}
	editedName := "ร้านที่ผู้ดูแลปรับจากข้อมูลตัวอย่าง"
	if _, err := s.UpdateRestaurant(ctx, ids[3], RestaurantInput{Name: &editedName}); err != nil {
		t.Fatal(err)
	}
	var adminCreatedID int64
	err = s.Pool.QueryRow(ctx, `INSERT INTO restaurants(name,province_id,category_id,address,description)
		SELECT 'ร้านตัวอย่างสมมติที่ผู้ดูแลเพิ่มเอง',p.id,c.id,'ที่อยู่ทดสอบ','ข้อมูลที่ผู้ดูแลเพิ่มเอง'
		FROM provinces p CROSS JOIN categories c WHERE p.name_th='เชียงใหม่' AND c.name_th='อาหารเหนือ' RETURNING id`).Scan(&adminCreatedID)
	if err != nil {
		t.Fatal(err)
	}
	// Edits to real source records must survive reseeding.
	var editedID int64
	err = s.Pool.QueryRow(ctx, `UPDATE restaurants SET description='คำอธิบายที่ผู้ดูแลแก้ไข'
		WHERE id=(SELECT min(id) FROM restaurants WHERE source_name IS DISTINCT FROM 'ThaiFood Guide author-created fictional sample') RETURNING id`).Scan(&editedID)
	if err != nil {
		t.Fatal(err)
	}
	seed, err := os.ReadFile(filepath.Join("..", "seed.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := s.Pool.Exec(ctx, string(seed)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.RestaurantByID(ctx, ids[0]); !errors.Is(err, ErrNotFound) {
		t.Fatalf("unused fictional restaurant was not retired: %v", err)
	}
	for _, id := range []int64{ids[1], ids[2], ids[3], adminCreatedID} {
		if _, err := s.RestaurantByID(ctx, id); err != nil {
			t.Fatalf("member-linked or admin-edited restaurant %d was removed: %v", id, err)
		}
	}
	r, err := s.RestaurantByID(ctx, editedID)
	if err != nil || r.Description != "คำอธิบายที่ผู้ดูแลแก้ไข" {
		t.Fatalf("reseed overwrote a restaurant edit: %v", err)
	}
	reviews, err := s.ListReviews(ctx, ids[1])
	if err != nil || len(reviews) != 1 || reviews[0].Rating != 5 {
		t.Fatalf("reseed changed existing reviews: %v", err)
	}
	favorites, err := s.ListFavorites(ctx, m.ID)
	if err != nil || len(favorites) != 1 || favorites[0].ID != ids[2] {
		t.Fatalf("reseed changed existing favorites: %v", err)
	}
}

func TestMemberUniquenessPasswordCASAndSafeJSON(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	m, err := s.CreateMember(ctx, "demo@example.com", "oldhash", "Demo")
	if err != nil {
		t.Fatal(err)
	}
	if m.Role != "member" {
		t.Fatalf("registration role: %s", m.Role)
	}
	encoded, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "hash") || strings.Contains(string(encoded), "Password") {
		t.Fatalf("member response contains password material: %s", encoded)
	}
	if _, err = s.CreateMember(ctx, "DEMO@example.com", "hash", "Other"); !errors.Is(err, ErrConflict) {
		t.Fatalf("duplicate: %v", err)
	}
	if err = s.UpdatePassword(ctx, m.ID, "wrong", "new"); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale hash: %v", err)
	}
	if err = s.UpdatePassword(ctx, m.ID, "oldhash", "newhash"); err != nil {
		t.Fatal(err)
	}
	m, err = s.MemberByEmail(ctx, "DEMO@example.com")
	if err != nil || m.PasswordHash != "newhash" {
		t.Fatalf("member lookup/CAS: %+v %v", m, err)
	}
}

func TestPartialRestaurantUpdateAndExplicitImageClear(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	list, _, err := s.ListRestaurants(ctx, RestaurantFilter{Page: 1, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	original := list[0]
	image := "https://example.com/owned-image.jpg"
	r, err := s.UpdateRestaurant(ctx, original.ID, RestaurantInput{ImageURL: &image, ImageURLSet: true})
	if err != nil {
		t.Fatal(err)
	}
	if r.Name != original.Name || r.ImageURL == nil || *r.ImageURL != image {
		t.Fatalf("partial image update: %+v", r)
	}
	name := "ชื่อใหม่สำหรับทดสอบ"
	r, err = s.UpdateRestaurant(ctx, original.ID, RestaurantInput{Name: &name})
	if err != nil {
		t.Fatal(err)
	}
	if r.Name != name || r.ImageURL == nil || *r.ImageURL != image {
		t.Fatalf("omitted image must be preserved: %+v", r)
	}
	r, err = s.UpdateRestaurant(ctx, original.ID, RestaurantInput{ImageURLSet: true})
	if err != nil {
		t.Fatal(err)
	}
	if r.ImageURL != nil || r.Address != original.Address {
		t.Fatalf("explicit null must clear only image: %+v", r)
	}
	items, total, err := s.ListRestaurants(ctx, RestaurantFilter{Province: r.Province, Category: r.Category, Query: name, Page: 1, Limit: 12})
	if err != nil || len(items) != 1 || total != 1 {
		t.Fatalf("combined Thai filters: %+v %d %v", items, total, err)
	}
	items, total, err = s.ListRestaurants(ctx, RestaurantFilter{Province: "ไม่มีจังหวัดนี้", Page: 1, Limit: 12})
	if err != nil || items == nil || len(items) != 0 || total != 0 {
		t.Fatalf("empty filter results: %+v %d %v", items, total, err)
	}
	if err = s.DeleteRestaurant(ctx, 999999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing: %v", err)
	}
}

func TestDatabaseConstraintsProtectInvalidDirectWrites(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	m, err := s.CreateMember(ctx, "constraints@example.com", "hash", "Constraints")
	if err != nil {
		t.Fatal(err)
	}
	list, _, err := s.ListRestaurants(ctx, RestaurantFilter{Page: 1, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		rating int
		text   string
	}{{0, "อาหารอร่อยและบริการดีมาก"}, {6, "อาหารอร่อยและบริการดีมาก"}, {5, "เก้า"}} {
		if _, err = s.CreateReview(ctx, list[0].ID, m.ID, tc.rating, tc.text); err == nil {
			t.Fatalf("invalid direct review accepted: %+v", tc)
		}
	}
	if err = s.CreateFavorite(ctx, m.ID, 999999); !errors.Is(err, ErrForeignKey) {
		t.Fatalf("favorite invalid restaurant: %v", err)
	}
	if _, err = s.Pool.Exec(ctx, `UPDATE members SET role='superadmin' WHERE id=$1`, m.ID); err == nil {
		t.Fatal("invalid role accepted")
	}
	if _, err = s.Pool.Exec(ctx, `UPDATE restaurants SET latitude=91,longitude=100 WHERE id=$1`, list[0].ID); err == nil {
		t.Fatal("out-of-range coordinate accepted")
	}
	if _, err = s.Pool.Exec(ctx, `UPDATE restaurants SET latitude=18,longitude=NULL WHERE id=$1`, list[0].ID); err == nil {
		t.Fatal("unpaired coordinate accepted")
	}
}

func TestRestaurantSearchReviewOwnershipAggregatesAndCascade(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	name, address, description := "ร้าน 100%_จริง", "ที่อยู่ทดสอบ", "รายละเอียดสำหรับทดสอบ"
	ps, _ := s.ListProvinces(ctx)
	cs, _ := s.ListCategories(ctx)
	in := RestaurantInput{Name: &name, ProvinceID: &ps[0].ID, CategoryID: &cs[0].ID, Address: &address, Description: &description}
	r, err := s.CreateRestaurant(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{"100%_", "' OR 1=1 --"} {
		items, total, err := s.ListRestaurants(ctx, RestaurantFilter{Query: q, Page: 1, Limit: 12})
		if err != nil {
			t.Fatal(err)
		}
		want := int64(1)
		if q[0] == '\'' {
			want = 0
		}
		if total != want || int64(len(items)) != want {
			t.Fatalf("literal search %q: total=%d items=%d", q, total, len(items))
		}
	}
	a, _ := s.CreateMember(ctx, "owner@example.com", "hash", "Owner")
	b, _ := s.CreateMember(ctx, "other@example.com", "hash", "Other")
	review, err := s.CreateReview(ctx, r.ID, a.ID, 5, "อาหารอร่อยและบริการดีมาก")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.CreateReview(ctx, r.ID, a.ID, 3, "รีวิวซ้ำต้องไม่ได้รับอนุญาต"); !errors.Is(err, ErrConflict) {
		t.Fatalf("review duplicate %v", err)
	}
	if err = s.DeleteReview(ctx, review.ID, b.ID, false); !errors.Is(err, ErrForbidden) {
		t.Fatalf("non-owner: %v", err)
	}
	if _, err = s.CreateReview(ctx, r.ID, b.ID, 3, "ร้านนี้บรรยากาศดีและอร่อย"); err != nil {
		t.Fatal(err)
	}
	r, err = s.RestaurantByID(ctx, r.ID)
	if err != nil || r.Rating == nil || *r.Rating != 4 || r.ReviewCount != 2 {
		t.Fatalf("aggregate: %+v %v", r, err)
	}
	if err = s.CreateFavorite(ctx, a.ID, r.ID); err != nil {
		t.Fatal(err)
	}
	if err = s.CreateFavorite(ctx, a.ID, r.ID); !errors.Is(err, ErrConflict) {
		t.Fatalf("favorite duplicate: %v", err)
	}
	if err = s.DeleteReview(ctx, review.ID, b.ID, true); err != nil {
		t.Fatal(err)
	}
	r, err = s.RestaurantByID(ctx, r.ID)
	if err != nil || *r.Rating != 3 || r.ReviewCount != 1 {
		t.Fatalf("aggregate after moderation: %+v %v", r, err)
	}
	if err = s.DeleteRestaurant(ctx, r.ID); err != nil {
		t.Fatal(err)
	}
	if _, err = s.ListReviews(ctx, r.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing restaurant reviews: %v", err)
	}
	f, err := s.ListFavorites(ctx, a.ID)
	if err != nil || f == nil || len(f) != 0 {
		t.Fatalf("cascade favorites: %+v %v", f, err)
	}
	if _, err = s.CreateReview(ctx, r.ID, a.ID, 4, "ไม่ควรรีวิวร้านที่ถูกลบแล้ว"); !errors.Is(err, ErrForeignKey) {
		t.Fatalf("FK: %v", err)
	}
}

func TestConcurrentDuplicateReviewOnlyOneSucceeds(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	m, _ := s.CreateMember(ctx, "race@example.com", "hash", "Race")
	list, _, _ := s.ListRestaurants(ctx, RestaurantFilter{Page: 1, Limit: 1})
	var wg sync.WaitGroup
	errs := make(chan error, 8)
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.CreateReview(ctx, list[0].ID, m.ID, 5, "อาหารอร่อยบริการดีสำหรับทดสอบ")
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	successes, conflicts := 0, 0
	for err := range errs {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrConflict) {
			conflicts++
		} else {
			t.Fatal(err)
		}
	}
	if successes != 1 || conflicts != 7 {
		t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
	}
}
