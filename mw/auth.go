package mw

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"thaifood-backend/auth"
	"thaifood-backend/store"
)

type memberContextKey struct{}

func WithMember(ctx context.Context, member store.Member) context.Context {
	return context.WithValue(ctx, memberContextKey{}, member)
}

func Member(ctx context.Context) (store.Member, bool) {
	m, ok := ctx.Value(memberContextKey{}).(store.Member)
	return m, ok && m.ID > 0
}

// RequireLogin reloads the member for every request. Tokens identify a member;
// current database data, rather than token claims, determines their role.
func RequireLogin(tickets *auth.Issuer, repository interface {
	MemberByID(context.Context, int64) (store.Member, error)
}, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || len(r.Header.Values("Authorization")) != 1 {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "A valid Bearer token is required")
			return
		}
		id, err := tickets.Parse(parts[1])
		if err != nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "A valid Bearer token is required")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), databaseTimeout)
		defer cancel()
		member, err := repository.MemberByID(ctx, id)
		if errors.Is(err, store.ErrNotFound) || (err == nil && member.ID != id) {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "A valid Bearer token is required")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
			return
		}
		if member.Role != "member" && member.Role != "admin" {
			writeError(w, http.StatusForbidden, "FORBIDDEN", "Your account does not have access")
			return
		}
		next(w, r.WithContext(WithMember(r.Context(), member)))
	}
}

func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		member, ok := Member(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "A valid Bearer token is required")
			return
		}
		if member.Role != role {
			writeError(w, http.StatusForbidden, "FORBIDDEN", "Your account does not have access")
			return
		}
		next(w, r)
	}
}
