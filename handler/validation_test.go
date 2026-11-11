package handler

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"thaifood-backend/auth"
	"thaifood-backend/store"
)

func TestRegisterNormalizesAndNeverExposesHash(t *testing.T) {
	s := &testStore{}
	h := (&Handler{Store: s, Tickets: &auth.Issuer{Secret: []byte(strings.Repeat("x", 32))}}).Routes()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(`{"email":"  TEST@Example.com ","password":"12345678","displayName":"  ผู้ใช้ทดสอบ  "}`))
	h.ServeHTTP(w, r)
	if w.Code != 201 || s.email != "test@example.com" || s.displayName != "ผู้ใช้ทดสอบ" || !auth.CheckPassword(s.hash, "12345678") {
		t.Fatalf("registration failed %d: %s", w.Code, w.Body.String())
	}
	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response) != 4 || response["role"] != "member" || strings.Contains(w.Body.String(), s.hash) {
		t.Fatalf("unsafe member response %s", w.Body.String())
	}
}

func TestReviewUsesUnicodeLengthAndRejectsFractions(t *testing.T) {
	issuer := &auth.Issuer{Secret: []byte(strings.Repeat("x", 32))}
	token, err := issuer.Issue(1)
	if err != nil {
		t.Fatal(err)
	}
	s := &testStore{member: store.Member{ID: 1, Role: "member", DisplayName: "Demo"}}
	h := (&Handler{Store: s, Tickets: issuer}).Routes()
	for _, tt := range []struct {
		body   string
		status int
	}{
		{`{"rating":5,"text":"อร่อย"}`, 400},
		{`{"rating":5.5,"text":"อาหารอร่อย บริการดี"}`, 400},
		{`{"rating":0,"text":"อาหารอร่อย บริการดี"}`, 400},
		{`{"rating":5,"text":"  อาหารอร่อย บริการดี  "}`, 201},
	} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/restaurants/10/reviews", strings.NewReader(tt.body))
		r.Header.Set("Authorization", "Bearer "+token)
		h.ServeHTTP(w, r)
		if w.Code != tt.status {
			t.Errorf("body %s: %d %s", tt.body, w.Code, w.Body.String())
		}
	}
	if s.rating != 5 || s.reviewText != "อาหารอร่อย บริการดี" {
		t.Fatalf("stored review %+v", s)
	}
}

func TestStoreErrorMapping(t *testing.T) {
	for _, tt := range []struct {
		err    error
		status int
		code   string
	}{{store.ErrNotFound, 404, "NOT_FOUND"}, {store.ErrConflict, 409, "CONFLICT"}, {store.ErrForeignKey, 400, "INVALID_REFERENCE"}, {store.ErrForbidden, 403, "FORBIDDEN"}, {errors.New("secret"), 500, "INTERNAL_ERROR"}} {
		w := httptest.NewRecorder()
		storeError(w, tt.err)
		if w.Code != tt.status || !strings.Contains(w.Body.String(), tt.code) || strings.Contains(w.Body.String(), "secret") {
			t.Fatalf("mapping %v: %d %s", tt.err, w.Code, w.Body.String())
		}
	}
}

func TestBodyRejectsNonJSONAndOversize(t *testing.T) {
	for _, tt := range []struct{ body, media string }{{strings.Repeat("x", (1<<20)+1), "application/json"}, {`{"displayName":"Demo"}`, "text/plain"}, {`{"displayName":"Demo"} null`, "application/json"}} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PATCH", "/", strings.NewReader(tt.body))
		r.Header.Set("Content-Type", tt.media)
		var dst struct {
			DisplayName string `json:"displayName"`
		}
		if err := body(w, r, &dst); err == nil {
			t.Errorf("accepted invalid body of %d bytes (%s)", len(tt.body), tt.media)
		}
	}
}

func TestRestaurantDescriptionMatchesDatabaseLimit(t *testing.T) {
	encoded, err := json.Marshal(map[string]string{"description": strings.Repeat("ก", 5001)})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("PUT", "/", strings.NewReader(string(encoded)))
	if _, err := restaurantInput(httptest.NewRecorder(), request, false); err == nil {
		t.Fatal("accepted description longer than database constraint")
	}
}

func TestRestaurantFiltersRejectInvalidUTF8(t *testing.T) {
	s := &testStore{}
	h := (&Handler{Store: s}).Routes()
	for _, field := range []string{"q", "province", "category"} {
		t.Run(field, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/restaurants?"+field+"=%FF", nil)
			h.ServeHTTP(w, r)
			if w.Code != 400 {
				t.Fatalf("invalid UTF-8 filter should be rejected before store call: status=%d body=%s", w.Code, w.Body.String())
			}
			if s.filter != (store.RestaurantFilter{}) {
				t.Fatal("invalid UTF-8 filter reached the store")
			}
			var response struct {
				Error struct{ Code, Message string }
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.Error.Code != "VALIDATION_ERROR" || response.Error.Message == "" {
				t.Fatalf("invalid error response: %s", w.Body.String())
			}
		})
	}
}
