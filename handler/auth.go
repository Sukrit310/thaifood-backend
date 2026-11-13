package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"thaifood-backend/auth"
	"thaifood-backend/mw"
	"thaifood-backend/store"
)

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && len(value) <= 254 && strings.Contains(value, "@")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}
	if err := body(w, r, &in); err != nil {
		bad(w, err)
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.DisplayName = strings.TrimSpace(in.DisplayName)
	if !validEmail(in.Email) || !passwordValid(in.Password) || !boundedText(in.DisplayName, 1, 100) {
		bad(w, errors.New("Valid email, password of 8–72 bytes and displayName of 1–100 characters are required"))
		return
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		storeError(w, err)
		return
	}
	member, err := h.Store.CreateMember(r.Context(), in.Email, hash, in.DisplayName)
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 201, member)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := body(w, r, &in); err != nil {
		bad(w, err)
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if !validEmail(in.Email) || len(in.Password) == 0 || len(in.Password) > 72 {
		bad(w, errors.New("Valid email and password are required"))
		return
	}
	member, err := h.Store.MemberByEmail(r.Context(), in.Email)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		storeError(w, err)
		return
	}
	// Compare a valid dummy hash when the email is absent to avoid a fast account-existence probe.
	hash := member.PasswordHash
	if errors.Is(err, store.ErrNotFound) {
		hash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	}
	matched := auth.CheckPassword(hash, in.Password)
	if err != nil || !matched {
		fail(w, 401, "INVALID_CREDENTIALS", "Email or password is incorrect")
		return
	}
	token, err := h.Tickets.Issue(member.ID)
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 200, map[string]any{"token": token, "member": member})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	member, _ := mw.Member(r.Context())
	write(w, 200, member)
}
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	var in struct {
		DisplayName string `json:"displayName"`
	}
	if err := body(w, r, &in); err != nil {
		bad(w, err)
		return
	}
	in.DisplayName = strings.TrimSpace(in.DisplayName)
	if !boundedText(in.DisplayName, 1, 100) {
		bad(w, errors.New("displayName must contain 1–100 characters"))
		return
	}
	member, _ := mw.Member(r.Context())
	updated, err := h.Store.UpdateMember(r.Context(), member.ID, in.DisplayName)
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 200, updated)
}
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := body(w, r, &in); err != nil {
		bad(w, err)
		return
	}
	if len(in.CurrentPassword) == 0 || len(in.CurrentPassword) > 72 || !passwordValid(in.NewPassword) {
		bad(w, errors.New("Current password is required and newPassword must contain 8–72 bytes"))
		return
	}
	member, _ := mw.Member(r.Context())
	if !auth.CheckPassword(member.PasswordHash, in.CurrentPassword) {
		fail(w, 401, "INVALID_CREDENTIALS", "Current password is incorrect")
		return
	}
	hash, err := auth.HashPassword(in.NewPassword)
	if err != nil {
		storeError(w, err)
		return
	}
	if err = h.Store.UpdatePassword(r.Context(), member.ID, member.PasswordHash, hash); err != nil {
		storeError(w, err)
		return
	}
	write(w, 204, nil)
}
func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	list, err := h.Store.ListMembers(r.Context())
	if err != nil {
		storeError(w, err)
		return
	}
	items(w, list)
}

func (h *Handler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	targetID, err := id(r, "id")
	if err != nil {
		bad(w, err)
		return
	}
	// A map also rejects differently cased keys, which struct decoding accepts.
	var in map[string]json.RawMessage
	if err := body(w, r, &in); err != nil {
		bad(w, err)
		return
	}
	var role string
	if len(in) != 1 || json.Unmarshal(in["role"], &role) != nil || (role != "member" && role != "admin") {
		bad(w, errors.New("Only role is allowed and must be member or admin"))
		return
	}
	actor, _ := mw.Member(r.Context())
	if actor.ID == targetID {
		storeError(w, store.ErrForbidden)
		return
	}
	updated, err := h.Store.UpdateMemberRole(r.Context(), actor.ID, targetID, role)
	if err != nil {
		storeError(w, err)
		return
	}
	write(w, 200, updated)
}
