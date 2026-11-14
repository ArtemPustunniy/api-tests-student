package favorites

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite adds one advertisement to the current user's favorites.
func HttpPostFavorite(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(models.AddFavoriteRequest{AdvertisementID: advertisementID}).
		Expect(t).
		End().Response
}

// HttpGetFavorites returns a large first page of the current user's favorites.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", "1").
		Query("page_size", strconv.Itoa(50)).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite removes one advertisement from the current user's favorites.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.FavoritePath(advertisementID)).
		Expect(t).
		End().Response
}
