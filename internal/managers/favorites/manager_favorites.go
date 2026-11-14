package favorites

import (
	"testing"

	client "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/common"
	"api-tests-template/internal/managers/favorites/models"
)

type Result[T any] struct {
	StatusCode int
	Body       T
}

// AddFavorite adds an advertisement and decodes a successful response.
func AddFavorite(t *testing.T, token, advertisementID string) Result[models.FavoriteItem] {
	t.Helper()
	response := client.HttpPostFavorite(t, token, advertisementID)
	result := Result[models.FavoriteItem]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// AddFavoriteError adds an invalid favorite and decodes the API error.
func AddFavoriteError(t *testing.T, token, advertisementID string) Result[common.ErrorResponse] {
	t.Helper()
	response := client.HttpPostFavorite(t, token, advertisementID)
	result := Result[common.ErrorResponse]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// GetFavorites returns the current user's favorites or an authorization error.
func GetFavorites(t *testing.T, token string) Result[models.FavoriteList] {
	t.Helper()
	response := client.HttpGetFavorites(t, token)
	result := Result[models.FavoriteList]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// GetFavoritesError decodes an error response from the favorites list endpoint.
func GetFavoritesError(t *testing.T, token string) Result[common.ErrorResponse] {
	t.Helper()
	response := client.HttpGetFavorites(t, token)
	result := Result[common.ErrorResponse]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// DeleteFavorite removes an advertisement and returns the response status.
func DeleteFavorite(t *testing.T, token, advertisementID string) int {
	t.Helper()
	response := client.HttpDeleteFavorite(t, token, advertisementID)
	defer response.Body.Close()
	return response.StatusCode
}

// DeleteFavoriteError removes an invalid favorite and decodes the API error.
func DeleteFavoriteError(t *testing.T, token, advertisementID string) Result[common.ErrorResponse] {
	t.Helper()
	response := client.HttpDeleteFavorite(t, token, advertisementID)
	result := Result[common.ErrorResponse]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}
