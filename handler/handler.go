package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"thaifood-backend/auth"
	"thaifood-backend/mw"
	"thaifood-backend/store"
)

type Handler struct {
	Store   store.Repository
	Tickets *auth.Issuer
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	routes := map[string]map[string]http.HandlerFunc{}
	add := func(method, path string, f http.HandlerFunc) {
		if routes[path] == nil {
			routes[path] = map[string]http.HandlerFunc{}
		}
		routes[path][method] = f
	}
	login := func(f http.HandlerFunc) http.HandlerFunc { return mw.RequireLogin(h.Tickets, h.Store, f) }
	admin := func(f http.HandlerFunc) http.HandlerFunc { return login(mw.RequireRole("admin", f)) }
	add("GET", "/health", h.Health)
	add("POST", "/api/auth/register", h.Register)
	add("POST", "/api/auth/login", h.Login)
	add("GET", "/api/provinces", h.ListProvinces)
	add("GET", "/api/categories", h.ListCategories)
	add("GET", "/api/restaurants", h.ListRestaurants)
	add("GET", "/api/restaurants/{id}", h.GetRestaurant)
	add("GET", "/api/restaurants/{id}/reviews", h.ListReviews)
	add("POST", "/api/restaurants/{id}/reviews", login(h.CreateReview))
	add("GET", "/api/me", login(h.Me))
	add("PATCH", "/api/me", login(h.UpdateMe))
	add("PUT", "/api/me/password", login(h.ChangePassword))
	add("DELETE", "/api/reviews/{id}", login(h.DeleteReview))
	add("GET", "/api/favorites", login(h.ListFavorites))
	add("POST", "/api/favorites", login(h.CreateFavorite))
	add("DELETE", "/api/favorites/{restaurantId}", login(h.DeleteFavorite))
	add("GET", "/api/admin/members", admin(h.ListMembers))
	add("POST", "/api/admin/restaurants", admin(h.CreateRestaurant))
	add("PUT", "/api/admin/restaurants/{id}", admin(h.UpdateRestaurant))
	add("DELETE", "/api/admin/restaurants/{id}", admin(h.DeleteRestaurant))
	add("DELETE", "/api/admin/reviews/{id}", admin(h.DeleteReview))
	for path, methods := range routes {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if f, ok := methods[r.Method]; ok {
				f(w, r)
				return
			}
			allowed := []string{}
			for _, m := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
				if _, ok := methods[m]; ok {
					allowed = append(allowed, m)
				}
			}
			w.Header().Set("Allow", strings.Join(allowed, ", "))
			fail(w, 405, "METHOD_NOT_ALLOWED", "Method is not allowed for this endpoint")
		})
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fail(w, 404, "NOT_FOUND", "Endpoint not found") })
	return mux
}

func write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}
func fail(w http.ResponseWriter, status int, code, message string) {
	write(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
func bad(w http.ResponseWriter, err error) { fail(w, 400, "VALIDATION_ERROR", err.Error()) }
func storeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		fail(w, 404, "NOT_FOUND", "Requested resource not found")
	case errors.Is(err, store.ErrConflict):
		fail(w, 409, "CONFLICT", "This record already exists or was changed by another request")
	case errors.Is(err, store.ErrForeignKey):
		fail(w, 400, "INVALID_REFERENCE", "Province or category does not exist")
	case errors.Is(err, store.ErrForbidden):
		fail(w, 403, "FORBIDDEN", "You do not have permission for this action")
	default:
		fail(w, 500, "INTERNAL_ERROR", "Unable to complete the request")
	}
}
func items[T any](w http.ResponseWriter, list []T) {
	if list == nil {
		list = []T{}
	}
	write(w, 200, map[string]any{"items": list})
}

func body(w http.ResponseWriter, r *http.Request, dst any) error {
	if value := r.Header.Get("Content-Type"); value != "" {
		media, _, err := mime.ParseMediaType(value)
		if err != nil || media != "application/json" {
			return errors.New("Content-Type must be application/json")
		}
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		return errors.New("Request body must be valid JSON no larger than 1 MiB")
	}
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 || raw[0] != '{' || !utf8.Valid(raw) {
		return errors.New("Request body must be a JSON object")
	}
	d := json.NewDecoder(strings.NewReader(string(raw)))
	d.DisallowUnknownFields()
	if err = d.Decode(dst); err != nil {
		return errors.New("Request contains malformed JSON, unknown fields or invalid field types")
	}
	if d.Decode(&struct{}{}) != io.EOF {
		return errors.New("Request must contain one JSON object")
	}
	return nil
}

func id(r *http.Request, key string) (int64, error) {
	n, err := strconv.ParseInt(r.PathValue(key), 10, 64)
	if err != nil || n <= 0 {
		return 0, errors.New("ID must be a positive integer")
	}
	return n, nil
}
func boundedText(value string, min, max int) bool {
	n := utf8.RuneCountInString(value)
	return utf8.ValidString(value) && n >= min && n <= max && !strings.ContainsRune(value, '\x00')
}
func passwordValid(value string) bool { return len(value) >= 8 && len(value) <= 72 }
func requiredText(value *string, name string, max int) error {
	if value == nil {
		return fmt.Errorf("%s is required", name)
	}
	*value = strings.TrimSpace(*value)
	if !boundedText(*value, 1, max) {
		return fmt.Errorf("%s must contain 1–%d characters", name, max)
	}
	return nil
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.Ping(r.Context()); err != nil {
		fail(w, 503, "UNAVAILABLE", "Database is unavailable")
		return
	}
	write(w, 200, map[string]string{"status": "ok"})
}
