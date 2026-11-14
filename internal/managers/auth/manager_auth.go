package auth

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// Register creates a user and returns the token and user data from the API.
func Register(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	t.Helper()

	registerRawResponse := auth.HttpPostAuthRegister(t, request)
	defer registerRawResponse.Body.Close()

	responseBody, err := io.ReadAll(registerRawResponse.Body)
	require.NoError(t, err, "failed to read registration response")

	registerResponse := models.LoginOkResponse{}
	require.NoError(t, json.Unmarshal(responseBody, &registerResponse), "failed to decode registration response")
	require.NotEmpty(t, registerResponse.Token, "registration response must contain a token")

	return registerResponse
}
