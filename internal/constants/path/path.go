package path

import "fmt"

const (
	AuthLoginPath        = "/auth/login"
	AuthRegisterPath     = "/auth/register"
	MyAdvertisementsPath = "/my/advertisements"
	AdvertisementsPath   = "/advertisements"
	AdvertisementPath    = "/advertisement"
	FavoritesPath        = "/favorites"
)

// AdvertisementPhotosPath returns the endpoint for an advertisement's photos.
func AdvertisementPhotosPath(advertisementID string) string {
	return fmt.Sprintf("/advertisements/%s/photos", advertisementID)
}

// FavoritePath returns the endpoint for one favorite advertisement.
func FavoritePath(advertisementID string) string {
	return fmt.Sprintf("/favorites/%s", advertisementID)
}
