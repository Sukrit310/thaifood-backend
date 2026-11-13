// ไฟล์ handler/restaurants.go: อ่านจังหวัด/ประเภทอาหาร ค้นหาร้าน และตรวจข้อมูลสำหรับให้ admin จัดการร้าน
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"thaifood-backend/store"
)

// ListProvinces ส่งรายการจังหวัดพร้อมชื่อไทยและภูมิภาค
func (h *Handler) ListProvinces(w http.ResponseWriter, r *http.Request) {
	list, err := h.Store.ListProvinces(r.Context())
	if err != nil {
		storeError(w, err)
		return
	}
	items(w, list)
}

// ListCategories ส่งรายการประเภทอาหารที่ใช้เป็นตัวกรองและอ้างอิงตอนสร้างร้าน
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	list, err := h.Store.ListCategories(r.Context())
	if err != nil {
		storeError(w, err)
		return
	}
	items(w, list)
}

// ListRestaurants ตรวจ query parameters ใช้ค่าเริ่มต้น page=1/limit=12 จำกัด limit ที่ 50 และส่งผลพร้อมจำนวนรวม
func (h *Handler) ListRestaurants(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		bad(w, errors.New("Malformed query string"))
		return
	}
	for key, values := range query {
		if len(values) != 1 {
			bad(w, errors.New("Query parameters must occur once"))
			return
		}
		switch key {
		case "page", "limit", "province", "category", "q":
		default:
			bad(w, errors.New("Unknown query parameter"))
			return
		}
	}
	filter := store.RestaurantFilter{Province: strings.TrimSpace(query.Get("province")), Category: strings.TrimSpace(query.Get("category")), Query: strings.TrimSpace(query.Get("q")), Page: 1, Limit: 12}
	for key, target := range map[string]*int{"page": &filter.Page, "limit": &filter.Limit} {
		if values, ok := query[key]; ok {
			n, e := strconv.Atoi(values[0])
			if e != nil || n <= 0 || n > 1000000000 {
				bad(w, errors.New("page and limit must be positive integers no greater than 1000000000"))
				return
			}
			*target = n
		}
	}
	if filter.Limit > 50 {
		filter.Limit = 50
	}
	if !boundedText(filter.Province, 0, 100) || !boundedText(filter.Category, 0, 100) || !boundedText(filter.Query, 0, 200) {
		bad(w, errors.New("Filter values are too long"))
		return
	}
	list, total, err := h.Store.ListRestaurants(r.Context(), filter)
	if err != nil {
		storeError(w, err)
		return
	}
	if list == nil {
		list = []store.Restaurant{}
	}
	write(w, 200, map[string]any{"items": list, "page": filter.Page, "limit": filter.Limit, "total": total})
}

// GetRestaurant ตรวจรหัสร้านและส่งรายละเอียดพร้อมคะแนนเฉลี่ยและจำนวนรีวิว
func (h *Handler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := id(r, "id")
	if err != nil {
		bad(w, err)
		return
	}
	restaurant, err := h.Store.RestaurantByID(r.Context(), restaurantID)
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 200, restaurant)
}

// restaurantInput ตรวจฟิลด์ร้านตามโหมดสร้าง/แก้ไข รวมถึงพิกัดที่ต้องส่งเป็นคู่และ URL รูปภาพ
func restaurantInput(w http.ResponseWriter, r *http.Request, create bool) (store.RestaurantInput, error) {
	var raw map[string]json.RawMessage
	if err := body(w, r, &raw); err != nil {
		return store.RestaurantInput{}, err
	}
	in := store.RestaurantInput{}
	if len(raw) == 0 {
		return in, errors.New("At least one restaurant field is required")
	}
	allowed := map[string]bool{"name": true, "provinceId": true, "categoryId": true, "address": true, "description": true, "latitude": true, "longitude": true, "imageUrl": true}
	for key, value := range raw {
		if !allowed[key] {
			return in, errors.New("Unknown restaurant field")
		}
		if string(value) == "null" && key != "imageUrl" {
			return in, errors.New(key + " cannot be null")
		}
	}
	encoded, _ := json.Marshal(raw)
	if err := json.Unmarshal(encoded, &in); err != nil {
		return in, errors.New("Restaurant field types are invalid")
	}
	for _, field := range []struct {
		key   string
		value *string
		max   int
	}{{"name", in.Name, 200}, {"address", in.Address, 1000}, {"description", in.Description, 5000}} {
		if create || field.value != nil {
			if err := requiredText(field.value, field.key, field.max); err != nil {
				return in, err
			}
		}
	}
	for _, field := range []struct {
		key   string
		value *int64
	}{{"provinceId", in.ProvinceID}, {"categoryId", in.CategoryID}} {
		if create && field.value == nil {
			return in, errors.New(field.key + " is required")
		}
		if field.value != nil && *field.value <= 0 {
			return in, errors.New(field.key + " must be positive")
		}
	}
	if (in.Latitude == nil) != (in.Longitude == nil) {
		return in, errors.New("latitude and longitude must be provided together")
	}
	if in.Latitude != nil && (*in.Latitude < -90 || *in.Latitude > 90 || *in.Longitude < -180 || *in.Longitude > 180) {
		return in, errors.New("Coordinates are out of range")
	}
	// แยก imageUrl ที่ไม่ได้ส่ง (คงรูปเดิม) ออกจากการส่ง null (ล้างรูป)
	_, in.ImageURLSet = raw["imageUrl"]
	if in.ImageURL != nil {
		*in.ImageURL = strings.TrimSpace(*in.ImageURL)
		u, err := url.Parse(*in.ImageURL)
		if err != nil || u.Host == "" || u.User != nil || (u.Scheme != "https" && u.Scheme != "http") || len(*in.ImageURL) > 2048 {
			return in, errors.New("imageUrl must be an HTTP or HTTPS URL no longer than 2048 bytes")
		}
	}
	return in, nil
}

// CreateRestaurant ตรวจฟิลด์บังคับแล้วสร้างร้านใหม่ผ่าน endpoint สำหรับ admin
func (h *Handler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	in, err := restaurantInput(w, r, true)
	if err != nil {
		bad(w, err)
		return
	}
	restaurant, err := h.Store.CreateRestaurant(r.Context(), in)
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 201, restaurant)
}

// UpdateRestaurant ตรวจข้อมูลที่ส่งมาและแก้เฉพาะฟิลด์ที่ระบุ โดยฟิลด์ที่ละไว้คงค่าเดิม
func (h *Handler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := id(r, "id")
	if err != nil {
		bad(w, err)
		return
	}
	in, err := restaurantInput(w, r, false)
	if err != nil {
		bad(w, err)
		return
	}
	restaurant, err := h.Store.UpdateRestaurant(r.Context(), restaurantID, in)
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 200, restaurant)
}

// DeleteRestaurant ลบร้านตาม ID โดยฐานข้อมูลจะลบรีวิวและรายการโปรดที่อ้างถึงตาม foreign key
func (h *Handler) DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := id(r, "id")
	if err != nil {
		bad(w, err)
		return
	}
	if err = h.Store.DeleteRestaurant(r.Context(), restaurantID); err != nil {
		storeError(w, err)
		return
	}
	write(w, 204, nil)
}
