package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

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
