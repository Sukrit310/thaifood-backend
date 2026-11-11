package store

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrForeignKey = errors.New("invalid reference")
	ErrForbidden  = errors.New("forbidden")
)

type Member struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}

type Province struct {
	ID     int64  `json:"id"`
	NameTh string `json:"nameTh"`
	Region string `json:"region"`
}

type Category struct {
	ID     int64  `json:"id"`
	NameTh string `json:"nameTh"`
}

type Restaurant struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Province    string   `json:"province"`
	Category    string   `json:"category"`
	Address     string   `json:"address"`
	Description string   `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	ImageURL    *string  `json:"imageUrl"`
	Rating      *float64 `json:"rating"`
	ReviewCount int64    `json:"reviewCount"`
}

type ReviewMember struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"displayName"`
}

type Review struct {
	ID           int64        `json:"id"`
	RestaurantID int64        `json:"restaurantId"`
	Rating       int          `json:"rating"`
	Text         string       `json:"text"`
	CreatedAt    time.Time    `json:"createdAt"`
	Member       ReviewMember `json:"member"`
}

type RestaurantFilter struct {
	Province, Category, Query string
	Page, Limit               int
}

// RestaurantInput uses pointers to preserve omitted fields on updates. ImageURLSet
// distinguishes an omitted imageUrl from an explicit null that clears the image.
type RestaurantInput struct {
	Name        *string  `json:"name"`
	ProvinceID  *int64   `json:"provinceId"`
	CategoryID  *int64   `json:"categoryId"`
	Address     *string  `json:"address"`
	Description *string  `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	ImageURL    *string  `json:"imageUrl"`
	ImageURLSet bool     `json:"-"`
}

type Repository interface {
	Ping(context.Context) error
	CreateMember(context.Context, string, string, string) (Member, error)
	MemberByEmail(context.Context, string) (Member, error)
	MemberByID(context.Context, int64) (Member, error)
	UpdateMember(context.Context, int64, string) (Member, error)
	UpdatePassword(context.Context, int64, string, string) error
	ListMembers(context.Context) ([]Member, error)
	ListProvinces(context.Context) ([]Province, error)
	ListCategories(context.Context) ([]Category, error)
	ListRestaurants(context.Context, RestaurantFilter) ([]Restaurant, int64, error)
	RestaurantByID(context.Context, int64) (Restaurant, error)
	CreateRestaurant(context.Context, RestaurantInput) (Restaurant, error)
	UpdateRestaurant(context.Context, int64, RestaurantInput) (Restaurant, error)
	DeleteRestaurant(context.Context, int64) error
	ListReviews(context.Context, int64) ([]Review, error)
	CreateReview(context.Context, int64, int64, int, string) (Review, error)
	DeleteReview(context.Context, int64, int64, bool) error
	ListFavorites(context.Context, int64) ([]Restaurant, error)
	CreateFavorite(context.Context, int64, int64) error
	DeleteFavorite(context.Context, int64, int64) error
}
