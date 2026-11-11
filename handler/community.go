package handler

import (
	"errors"
	"net/http"
	"strings"

	"thaifood-backend/mw"
	"thaifood-backend/store"
)

func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := id(r, "id")
	if err != nil {
		bad(w, err)
		return
	}
	list, err := h.Store.ListReviews(r.Context(), restaurantID)
	if err != nil {
		storeError(w, err)
		return
	}
	items(w, list)
}
func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := id(r, "id")
	if err != nil {
		bad(w, err)
		return
	}
	var in struct {
		Rating int    `json:"rating"`
		Text   string `json:"text"`
	}
	if err = body(w, r, &in); err != nil {
		bad(w, err)
		return
	}
	in.Text = strings.TrimSpace(in.Text)
	if in.Rating < 1 || in.Rating > 5 || !boundedText(in.Text, 10, 5000) {
		bad(w, errors.New("rating must be an integer from 1–5 and text must contain 10–5000 characters"))
		return
	}
	member, _ := mw.Member(r.Context())
	review, err := h.Store.CreateReview(r.Context(), restaurantID, member.ID, in.Rating, in.Text)
	if errors.Is(err, store.ErrForeignKey) {
		err = store.ErrNotFound
	}
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 201, review)
}
func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := id(r, "id")
	if err != nil {
		bad(w, err)
		return
	}
	member, _ := mw.Member(r.Context())
	if err = h.Store.DeleteReview(r.Context(), reviewID, member.ID, member.Role == "admin"); err != nil {
		storeError(w, err)
		return
	}
	write(w, 204, nil)
}
func (h *Handler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	member, _ := mw.Member(r.Context())
	list, err := h.Store.ListFavorites(r.Context(), member.ID)
	if err != nil {
		storeError(w, err)
		return
	}
	items(w, list)
}
func (h *Handler) CreateFavorite(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RestaurantID int64 `json:"restaurantId"`
	}
	if err := body(w, r, &in); err != nil {
		bad(w, err)
		return
	}
	if in.RestaurantID <= 0 {
		bad(w, errors.New("restaurantId must be a positive integer"))
		return
	}
	member, _ := mw.Member(r.Context())
	err := h.Store.CreateFavorite(r.Context(), member.ID, in.RestaurantID)
	if errors.Is(err, store.ErrForeignKey) {
		err = store.ErrNotFound
	}
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 201, in)
}
func (h *Handler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := id(r, "restaurantId")
	if err != nil {
		bad(w, err)
		return
	}
	member, _ := mw.Member(r.Context())
	if err = h.Store.DeleteFavorite(r.Context(), member.ID, restaurantID); err != nil {
		storeError(w, err)
		return
	}
	write(w, 204, nil)
}
