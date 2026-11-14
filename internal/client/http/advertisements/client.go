package advertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements выполняет GET-запрос поиска объявлений по вхождению подстроки в название, ожидает 200
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement creates an advertisement with all supplied fields and photos.
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	t.Helper()

	apiRequest := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		MultipartFormData("title", request.Title).
		MultipartFormData("description", request.Description)
	if request.Price != "" {
		apiRequest.MultipartFormData("price", request.Price)
	}
	if request.Quantity != "" {
		apiRequest.MultipartFormData("quantity", request.Quantity)
	}
	if len(request.PhotoPaths) > 0 {
		apiRequest.MultipartFile("photos", request.PhotoPaths...)
	}

	return apiRequest.Expect(t).End().Response
}

// HttpGetAdvertisement gets one advertisement by ID.
func HttpGetAdvertisement(t *testing.T, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", advertisementID).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement deletes one advertisement by ID.
func HttpDeleteAdvertisement(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", advertisementID).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos returns photos uploaded for an advertisement.
func HttpGetAdvertisementPhotos(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPhotosPath(advertisementID)).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementsWithLimit searches advertisements and applies a deterministic page limit.
func HttpGetAdvertisementsWithLimit(t *testing.T, textForSearch string, limit int) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Query("limit", strconv.Itoa(limit)).
		Expect(t).
		End().Response
}
