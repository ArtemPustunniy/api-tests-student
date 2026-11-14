package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementManager "api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favoriteManager "api-tests-template/internal/managers/favorites"
	favoriteModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/helpers"
)

const nonexistentAdvertisementID = "11111111-1111-4111-8111-111111111111"

type testUsers struct {
	ownerToken    string
	favoriteToken string
}

func TestFavoriteAddAndRemoveCycle(t *testing.T) {
	users := prepareUsers(t)
	advertisement := createAdvertisement(t, users.ownerToken)

	base.Precondition("Add another user's advertisement to favorites")
	base.Step("Add the advertisement and verify it in GET /favorites")
	added := favoriteManager.AddFavorite(t, users.favoriteToken, advertisement.Id)
	require.Equal(t, http.StatusCreated, added.StatusCode, "adding a favorite must return 201")
	require.Equal(t, advertisement.Id, added.Body.Advertisement.Id, "API added a different advertisement")
	t.Cleanup(func() {
		favoriteManager.DeleteFavorite(t, users.favoriteToken, advertisement.Id)
	})

	list := favoriteManager.GetFavorites(t, users.favoriteToken)
	require.Equal(t, http.StatusOK, list.StatusCode, "favorites list must return 200")
	require.True(t, containsFavorite(list.Body.Items, advertisement.Id), "added advertisement is absent from favorites")

	status := favoriteManager.DeleteFavorite(t, users.favoriteToken, advertisement.Id)
	require.Equal(t, http.StatusNoContent, status, "removing a favorite must return 204")

	list = favoriteManager.GetFavorites(t, users.favoriteToken)
	require.Equal(t, http.StatusOK, list.StatusCode, "favorites list must return 200")
	require.False(t, containsFavorite(list.Body.Items, advertisement.Id), "removed advertisement is still in favorites")
	base.Result("Advertisement appeared after adding and disappeared after removing")
}

func TestDeletedAdvertisementDisappearsFromFavorites(t *testing.T) {
	users := prepareUsers(t)
	advertisement := createAdvertisement(t, users.ownerToken)
	base.Precondition("Another user has the owner's advertisement in favorites")

	added := favoriteManager.AddFavorite(t, users.favoriteToken, advertisement.Id)
	require.Equal(t, http.StatusCreated, added.StatusCode, "adding a favorite must return 201")

	base.Step("Delete the advertisement as its owner")
	deleteStatus := advertisementManager.DeleteAdvertisement(t, users.ownerToken, advertisement.Id)
	require.Equal(t, http.StatusNoContent, deleteStatus, "advertisement deletion must return 204")

	list := favoriteManager.GetFavorites(t, users.favoriteToken)
	require.Equal(t, http.StatusOK, list.StatusCode, "favorites list must return 200")
	require.False(t, containsFavorite(list.Body.Items, advertisement.Id), "deleted advertisement must disappear from favorites")
	base.Result("Deleted advertisement disappeared from the other user's favorites")
}

func TestFavoriteWithNonexistentAdvertisement(t *testing.T) {
	users := prepareUsers(t)
	base.Precondition("Use a valid UUID that does not belong to any advertisement")

	t.Run("add", func(t *testing.T) {
		base.Step("Try to add a nonexistent advertisement to favorites")
		result := favoriteManager.AddFavoriteError(t, users.favoriteToken, nonexistentAdvertisementID)
		t.Logf("add nonexistent response: status=%d error=%q message=%q", result.StatusCode, result.Body.Error, result.Body.Message)
		require.Equal(t, http.StatusNotFound, result.StatusCode, "adding a nonexistent advertisement must return 404")
		require.Equal(t, "not_found", result.Body.Error, "unexpected error code")
		require.Equal(t, "Advertisement not found", result.Body.Message, "unexpected error message")
	})

	t.Run("remove", func(t *testing.T) {
		base.Step("Try to remove a nonexistent advertisement from favorites")
		result := favoriteManager.DeleteFavoriteError(t, users.favoriteToken, nonexistentAdvertisementID)
		t.Logf("remove nonexistent response: status=%d error=%q message=%q", result.StatusCode, result.Body.Error, result.Body.Message)
		require.Equal(t, http.StatusNotFound, result.StatusCode, "removing a nonexistent favorite must return 404")
		require.Equal(t, "not_found", result.Body.Error, "unexpected error code")
		require.Equal(t, "Favorite entry not found", result.Body.Message, "unexpected error message")
	})
}

func TestFavoritesWithoutAuthorization(t *testing.T) {
	base.SetupSuite()
	base.Precondition("Do not provide an Authorization header")

	t.Run("list", func(t *testing.T) {
		base.Step("Request GET /favorites without a token")
		result := favoriteManager.GetFavoritesError(t, "")
		assertUnauthorized(t, result.StatusCode, result.Body.Error, result.Body.Message)
	})

	t.Run("add", func(t *testing.T) {
		base.Step("Request POST /favorites without a token")
		result := favoriteManager.AddFavoriteError(t, "", nonexistentAdvertisementID)
		assertUnauthorized(t, result.StatusCode, result.Body.Error, result.Body.Message)
	})

	t.Run("remove", func(t *testing.T) {
		base.Step("Request DELETE /favorites/{ad_id} without a token")
		result := favoriteManager.DeleteFavoriteError(t, "", nonexistentAdvertisementID)
		assertUnauthorized(t, result.StatusCode, result.Body.Error, result.Body.Message)
	})
}

func TestFavoriteBusinessRestrictions(t *testing.T) {
	users := prepareUsers(t)
	advertisement := createAdvertisement(t, users.ownerToken)
	base.Precondition("Create an advertisement owned by the configured user")

	t.Run("duplicate favorite", func(t *testing.T) {
		base.Step("Add the same advertisement to favorites twice")
		first := favoriteManager.AddFavorite(t, users.favoriteToken, advertisement.Id)
		require.Equal(t, http.StatusCreated, first.StatusCode, "first addition must return 201")
		t.Cleanup(func() {
			favoriteManager.DeleteFavorite(t, users.favoriteToken, advertisement.Id)
		})

		duplicate := favoriteManager.AddFavoriteError(t, users.favoriteToken, advertisement.Id)
		t.Logf("duplicate response: status=%d error=%q message=%q", duplicate.StatusCode, duplicate.Body.Error, duplicate.Body.Message)
		require.Equal(t, http.StatusConflict, duplicate.StatusCode, "duplicate favorite must return 409")
		require.Equal(t, "already_exists", duplicate.Body.Error, "unexpected error code")
		require.Equal(t, "Advertisement is already in favorites", duplicate.Body.Message, "unexpected error message")
	})

	t.Run("own advertisement", func(t *testing.T) {
		base.Step("Try to add the owner's own advertisement to favorites")
		result := favoriteManager.AddFavoriteError(t, users.ownerToken, advertisement.Id)
		t.Logf("own advertisement response: status=%d error=%q message=%q", result.StatusCode, result.Body.Error, result.Body.Message)
		require.Equal(t, http.StatusForbidden, result.StatusCode, "adding an own advertisement must return 403")
		require.Equal(t, "forbidden", result.Body.Error, "unexpected error code")
		require.Equal(t, "Cannot add own advertisement to favorites", result.Body.Message, "unexpected error message")
	})
}

// prepareUsers logs in the configured owner and registers an isolated favorites user.
func prepareUsers(t *testing.T) testUsers {
	t.Helper()
	base.SetupSuite()
	owner := auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	suffix := utils.RandomString(16)
	favoriteUser := auth.Register(t, authModels.RegisterRequest{
		Email:     "lab3-" + suffix + "@example.com",
		Password:  "Lab3-password-" + suffix,
		FirstName: "LabTest",
		LastName:  "FavoriteUser",
	})
	return testUsers{ownerToken: owner.Token, favoriteToken: favoriteUser.Token}
}

// createAdvertisement creates a unique owner advertisement and registers cleanup.
func createAdvertisement(t *testing.T, token string) advertisementModels.Advertisement {
	t.Helper()
	request := advertisementModels.CreateAdvertisementRequest{
		Title:       "lab3-favorite-" + utils.RandomString(12),
		Description: "Advertisement created for an atomic favorites scenario",
		Price:       "5000",
		Quantity:    "2",
		PhotoPaths:  []string{helpers.CreateTestPNG(t)},
	}
	result := advertisementManager.CreateAdvertisement(t, token, request)
	require.Equal(t, http.StatusCreated, result.StatusCode, "test advertisement creation must return 201")
	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, token, result.Body.Id)
	})
	return result.Body
}

// containsFavorite reports whether the list contains an advertisement ID.
func containsFavorite(items []favoriteModels.FavoriteItem, advertisementID string) bool {
	for _, item := range items {
		if item.Advertisement.Id == advertisementID {
			return true
		}
	}
	return false
}

// assertUnauthorized checks the common contract of protected favorites endpoints.
func assertUnauthorized(t *testing.T, statusCode int, errorCode, message string) {
	t.Helper()
	require.Equal(t, http.StatusUnauthorized, statusCode, "request without a token must return 401")
	require.Equal(t, "unauthorized", errorCode, "unexpected authorization error code")
	require.NotEmpty(t, message, "authorization error must contain a message")
}
