package models

import advertisementModels "api-tests-template/internal/managers/advertisements/models"

type FavoriteItem struct {
	FavoriteID    string                            `json:"favorite_id"`
	AddedAt       string                            `json:"added_at"`
	Advertisement advertisementModels.Advertisement `json:"advertisement"`
}

type FavoriteList struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}
