package models

import authModels "api-tests-template/internal/managers/auth/models"

// AdvertisementsList ответ GET /advertisements: страница объявлений с курсорной пагинацией
type AdvertisementsList struct {
	Items      []Advertisement `json:"items"`
	HasNext    bool            `json:"has_next"`
	NextCursor string          `json:"next_cursor,omitempty"`
	Total      int             `json:"total"`
}

type Advertisement struct {
	Id          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	UserId      string          `json:"user_id"`
	User        authModels.User `json:"user"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Photos      []Photo         `json:"photos"`
}

type Photo struct {
	Id              string `json:"id"`
	AdvertisementId string `json:"advertisement_id"`
	Url             string `json:"url"`
	SortOrder       int    `json:"sort_order"`
	CreatedAt       string `json:"created_at"`
}
