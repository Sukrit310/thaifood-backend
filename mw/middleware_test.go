package mw

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"thaifood-backend/auth"
	"thaifood-backend/store"
)

type memberRepository struct {
	member   store.Member
	err      error
	calledID int64
}

func (r *memberRepository) MemberByID(_ context.Context, id int64) (store.Member, error) {
	r.calledID = id
	return r.member, r.err
}

func assertError(t *testing.T, w *httptest.ResponseRecorder, status int) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Error struct{ Code, Message string }
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Error.Code == "" || body.Error.Message == "" {
		t.Fatalf("bad error envelope: %s", w.Body.String())
	}
}

func TestAuthenticationAndCurrentDatabaseRole(t *testing.T) {
	i := &auth.Issuer{Secret: []byte(strings.Repeat("s", 32))}
	token, err := i.Issue(7)
	if err != nil {
		t.Fatal(err)
	}
	repo := &memberRepository{member: store.Member{ID: 7, Role: "member"}}
	admin := RequireLogin(i, repo, RequireRole("admin", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	r := httptest.NewRequest("GET", "/api/admin/members", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	admin(w, r)
	assertError(t, w, 403)
	if repo.calledID != 7 {
		t.Fatal("did not reload database identity")
	}
	repo.member.Role = "admin"
	w = httptest.NewRecorder()
	admin(w, r)
	if w.Code != 204 {
		t.Fatalf("fresh admin role status=%d", w.Code)
	}
	repo.member.Role = "member"
	w = httptest.NewRecorder()
	admin(w, r)
	assertError(t, w, 403)
	repo.err = store.ErrNotFound
	w = httptest.NewRecorder()
	admin(w, r)
	assertError(t, w, 401)
	repo.err = errors.New("postgres://secret:password@db")
	w = httptest.NewRecorder()
	admin(w, r)
	assertError(t, w, 500)
	if strings.Contains(w.Body.String(), "password") {
		t.Fatal("database details leaked")
	}
}

func TestRejectMissingOrMalformedAuthorization(t *testing.T) {
	i := &auth.Issuer{Secret: []byte(strings.Repeat("s", 32))}
	token, _ := i.Issue(1)
	for _, header := range []string{"", "Basic abc", "Bearer", "Bearer ", "Bearer garbage", "Bearer " + token + " extra"} {
		t.Run(header[:min(len(header), 15)], func(t *testing.T) {
			repo := &memberRepository{member: store.Member{ID: 1, Role: "member"}}
			h := RequireLogin(i, repo, func(w http.ResponseWriter, r *http.Request) { t.Fatal("unauthenticated request reached handler") })
			r := httptest.NewRequest("GET", "/api/me", nil)
			r.Header.Set("Authorization", header)
			w := httptest.NewRecorder()
			h(w, r)
			assertError(t, w, 401)
			if repo.calledID != 0 {
				t.Fatal("DB called for invalid token")
			}
		})
	}
}

func TestRoleRequiresMemberContext(t *testing.T) {
	w := httptest.NewRecorder()
	RequireRole("admin", func(http.ResponseWriter, *http.Request) { t.Fatal("missing context allowed") })(w, httptest.NewRequest("GET", "/", nil))
	assertError(t, w, 401)
}

func TestCORS(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := CORS([]string{"http://localhost:3000", "https://frontend.vercel.app"}, next)
	for _, tc := range []struct {
		name, method, origin, requestedMethod, requestedHeaders string
		status                                                  int
	}{
		{"nonbrowser", "GET", "", "", "", 200},
		{"allowed", "GET", "http://localhost:3000", "", "", 200},
		{"preflight", "OPTIONS", "https://frontend.vercel.app", "PATCH", "authorization,content-type", 204},
		{"untrusted", "GET", "https://attacker.example", "", "", 403},
		{"suffix", "GET", "https://frontend.vercel.app.attacker.example", "", "", 403},
		{"bad method", "OPTIONS", "http://localhost:3000", "TRACE", "", 403},
		{"bad header", "OPTIONS", "http://localhost:3000", "POST", "x-unapproved", 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/api/me", nil)
			r.Header.Set("Origin", tc.origin)
			r.Header.Set("Access-Control-Request-Method", tc.requestedMethod)
			r.Header.Set("Access-Control-Request-Headers", tc.requestedHeaders)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status < 400 && tc.origin != "" && w.Header().Get("Access-Control-Allow-Origin") != tc.origin {
				t.Fatal("missing exact allowed origin")
			}
			if tc.status == 403 {
				assertError(t, w, 403)
				if w.Header().Get("Access-Control-Allow-Origin") != "" {
					t.Fatal("rejected origin reflected")
				}
			}
		})
	}
}

func TestRecoverHidesPanic(t *testing.T) {
	w := httptest.NewRecorder()
	Recover(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("database secret") })).ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	assertError(t, w, 500)
	if strings.Contains(w.Body.String(), "secret") {
		t.Fatal("panic details leaked")
	}
}

func TestLoggingOmitsCredentialsAndQuery(t *testing.T) {
	var output bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&output)
	defer log.SetOutput(previous)
	h := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	r := httptest.NewRequest("POST", "/api/auth/login?password=query-secret", strings.NewReader(`{"password":"body-secret"}`))
	r.Header.Set("Authorization", "Bearer token-secret")
	h.ServeHTTP(httptest.NewRecorder(), r)
	line := output.String()
	for _, secret := range []string{"query-secret", "body-secret", "token-secret"} {
		if strings.Contains(line, secret) {
			t.Fatalf("log leaked %s", secret)
		}
	}
	if !strings.Contains(line, "status=201") || !strings.Contains(line, "/api/auth/login") {
		t.Fatalf("missing request metadata: %s", line)
	}
}
