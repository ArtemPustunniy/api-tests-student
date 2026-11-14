package advertisementCreation

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementManager "api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/helpers"
)

const (
	testDescription = "Atomic API test advertisement with every supported field"
	testPrice       = "12500.50"
	testQuantity    = "3"
)

func TestCreateAdvertisementWithAllFields(t *testing.T) {
	base.SetupSuite()
	login := auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	title := "lab3-item-" + utils.RandomString(12)
	request := validAdvertisement(t, title)

	base.Precondition("Create a unique advertisement with title, description, price, quantity and photo")
	base.Step("Send POST /advertisement with all supported fields")
	created := advertisementManager.CreateAdvertisement(t, login.Token, request)
	require.Equal(t, http.StatusCreated, created.StatusCode, "advertisement creation must return 201")
	require.NotEmpty(t, created.Body.Id, "created advertisement must have an ID")
	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, login.Token, created.Body.Id)
	})

	t.Run("advertisement is available by ID with the same fields", func(t *testing.T) {
		base.Step("Request the created advertisement by ID")
		actual := advertisementManager.GetAdvertisement(t, created.Body.Id)
		require.Equal(t, http.StatusOK, actual.StatusCode, "GET advertisement must return 200")
		assertAdvertisementFields(t, actual.Body, request)
	})

	t.Run("advertisement is found by its full title", func(t *testing.T) {
		base.Step("Search for the advertisement by the full title")
		searchResult := advertisementManager.SearchAdvertisements(t, title)
		require.Equal(t, http.StatusOK, searchResult.StatusCode, "search must return 200")

		for _, item := range searchResult.Body.Items {
			if item.Id == created.Body.Id {
				assertAdvertisementFields(t, item, request)
				return
			}
		}
		require.Failf(t, "created advertisement is absent from search", "ID %q was not found by full title %q", created.Body.Id, title)
	})

	t.Run("uploaded photo is stored on the server", func(t *testing.T) {
		base.Step("Request photos of the created advertisement")
		photos := advertisementManager.GetAdvertisementPhotos(t, login.Token, created.Body.Id)
		require.Equal(t, http.StatusOK, photos.StatusCode, "photo endpoint must return 200")
		require.Len(t, photos.Body, 1, "exactly one uploaded photo is expected")
		require.NotEmpty(t, photos.Body[0].Id, "photo must have an ID")
		require.NotEmpty(t, photos.Body[0].Url, "photo must have an accessible URL")
	})
	base.Result("Created advertisement, search result and uploaded photo match the request")
}

func TestCreateAdvertisementValidation(t *testing.T) {
	base.SetupSuite()
	login := auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	base.Precondition("Authorize and prepare isolated invalid advertisement requests")

	testCases := []struct {
		name            string
		expectedMessage string
		request         func(*testing.T) advertisementModels.CreateAdvertisementRequest
	}{
		{
			name:            "empty title",
			expectedMessage: "title is required",
			request: func(t *testing.T) advertisementModels.CreateAdvertisementRequest {
				return validAdvertisement(t, "")
			},
		},
		{
			name:            "missing description",
			expectedMessage: "description is required",
			request: func(t *testing.T) advertisementModels.CreateAdvertisementRequest {
				request := validAdvertisement(t, "lab3-no-description-"+utils.RandomString(8))
				request.Description = ""
				return request
			},
		},
		{
			name:            "negative price",
			expectedMessage: "price must be between 0 and 1000000 inclusive",
			request: func(t *testing.T) advertisementModels.CreateAdvertisementRequest {
				request := validAdvertisement(t, "lab3-negative-price-"+utils.RandomString(8))
				request.Price = "-1"
				return request
			},
		},
		{
			name:            "title longer than 50 characters",
			expectedMessage: "title must not exceed 50 characters provided 51",
			request: func(t *testing.T) advertisementModels.CreateAdvertisementRequest {
				return validAdvertisement(t, strings.Repeat("a", 51))
			},
		},
		{
			name:            "missing required photo",
			expectedMessage: "at least 1 photo is required",
			request: func(t *testing.T) advertisementModels.CreateAdvertisementRequest {
				request := validAdvertisement(t, "lab3-no-photo-"+utils.RandomString(8))
				request.PhotoPaths = nil
				return request
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			base.Step("Send invalid POST /advertisement: " + testCase.name)
			result := advertisementManager.CreateAdvertisementError(t, login.Token, testCase.request(t))
			t.Logf("validation response: status=%d error=%q message=%q details=%v", result.StatusCode, result.Body.Error, result.Body.Message, result.Body.Details)
			require.Equal(t, http.StatusBadRequest, result.StatusCode, "invalid advertisement must return 400")
			require.Equal(t, "invalid_request", result.Body.Error, "unexpected validation error code")
			require.Equal(t, testCase.expectedMessage, result.Body.Message, "unexpected validation error message")
			base.Result("API returned the expected validation error")
		})
	}
}

func TestCreateAdvertisementWithoutAuthorization(t *testing.T) {
	base.SetupSuite()
	request := validAdvertisement(t, "lab3-unauthorized-"+utils.RandomString(12))
	base.Precondition("Prepare a valid advertisement request without valid authorization")

	for _, testCase := range []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{name: "missing token", token: "", expectedMessage: "Authorization header required"},
		{name: "invalid token", token: "invalid-token", expectedMessage: "Invalid or expired token"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			base.Step("Send POST /advertisement: " + testCase.name)
			result := advertisementManager.CreateAdvertisementError(t, testCase.token, request)
			t.Logf("authorization response: status=%d error=%q message=%q", result.StatusCode, result.Body.Error, result.Body.Message)
			require.Equal(t, http.StatusUnauthorized, result.StatusCode, "unauthorized create request must return 401")
			require.Equal(t, "unauthorized", result.Body.Error, "unexpected authorization error code")
			require.Equal(t, testCase.expectedMessage, result.Body.Message, "unexpected authorization error message")
			base.Result("API rejected the request with the expected 401 response")
		})
	}
}

// validAdvertisement builds a valid request with every supported field.
func validAdvertisement(t *testing.T, title string) advertisementModels.CreateAdvertisementRequest {
	t.Helper()
	return advertisementModels.CreateAdvertisementRequest{
		Title:       title,
		Description: testDescription,
		Price:       testPrice,
		Quantity:    testQuantity,
		PhotoPaths:  []string{helpers.CreateTestPNG(t)},
	}
}

// assertAdvertisementFields checks that the API preserved all submitted fields.
func assertAdvertisementFields(t *testing.T, actual advertisementModels.Advertisement, expected advertisementModels.CreateAdvertisementRequest) {
	t.Helper()
	require.Equal(t, expected.Title, actual.Title, "title differs from the submitted value")
	require.Equal(t, expected.Description, actual.Description, "description differs from the submitted value")
	require.InDelta(t, 12500.50, actual.Price, 0.001, "price differs from the submitted value")
	require.Equal(t, 3, actual.Quantity, "quantity differs from the submitted value")
}
