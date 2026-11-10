package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет, что среди объявлений есть объявление с точно таким названием
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	titles := make([]string, 0, len(items))
	for _, item := range items {
		if item.Title == title {
			return
		}
		titles = append(titles, item.Title)
	}

	require.Failf(t, "Объявление не найдено",
		"Объявления с названием %q нет среди найденных (%d шт.): %q", title, len(items), titles)
}
