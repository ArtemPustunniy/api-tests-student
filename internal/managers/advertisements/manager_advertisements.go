package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/common"
)

type Result[T any] struct {
	StatusCode int
	Body       T
}

// GetAdvertisementsList ищет объявления по вхождению подстроки в название и возвращает типизированный список.
// Статус-код 200 проверяет клиент
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	advertisementsRawResponse := advertisements.HttpGetAdvertisements(t, token, text)
	defer advertisementsRawResponse.Body.Close()

	body, errResponse := io.ReadAll(advertisementsRawResponse.Body)
	require.NoError(t, errResponse)

	advertisementsList := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &advertisementsList))

	return advertisementsList
}

// CreateAdvertisement sends a create request and decodes the successful response.
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) Result[models.Advertisement] {
	t.Helper()
	response := advertisements.HttpPostAdvertisement(t, token, request)
	result := Result[models.Advertisement]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// CreateAdvertisementError sends an invalid create request and decodes the API error.
func CreateAdvertisementError(t *testing.T, token string, request models.CreateAdvertisementRequest) Result[common.ErrorResponse] {
	t.Helper()
	response := advertisements.HttpPostAdvertisement(t, token, request)
	result := Result[common.ErrorResponse]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// GetAdvertisement gets and decodes one advertisement.
func GetAdvertisement(t *testing.T, advertisementID string) Result[models.Advertisement] {
	t.Helper()
	response := advertisements.HttpGetAdvertisement(t, advertisementID)
	result := Result[models.Advertisement]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// DeleteAdvertisement deletes one advertisement and returns its status code.
func DeleteAdvertisement(t *testing.T, token, advertisementID string) int {
	t.Helper()
	response := advertisements.HttpDeleteAdvertisement(t, token, advertisementID)
	defer response.Body.Close()
	return response.StatusCode
}

// GetAdvertisementPhotos gets and decodes all photos for an advertisement.
func GetAdvertisementPhotos(t *testing.T, token, advertisementID string) Result[[]models.Photo] {
	t.Helper()
	response := advertisements.HttpGetAdvertisementPhotos(t, token, advertisementID)
	result := Result[[]models.Photo]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}

// SearchAdvertisements searches advertisements using a bounded result page.
func SearchAdvertisements(t *testing.T, title string) Result[models.AdvertisementsList] {
	t.Helper()
	response := advertisements.HttpGetAdvertisementsWithLimit(t, title, 50)
	result := Result[models.AdvertisementsList]{StatusCode: response.StatusCode}
	common.DecodeResponse(t, response, &result.Body)
	return result
}
