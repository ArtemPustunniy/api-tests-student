package api_runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/steinfletcher/apitest"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"
	maxAttempts      = 5
)

type retryTransport struct {
	base http.RoundTripper
}

type ApiTest struct {
	token     string
	userAgent string
	host      string
}

func New() *ApiTest {
	return &ApiTest{
		userAgent: defaultUserAgent,
		host:      os.Getenv("API_URL"),
	}
}

func GetRunner() *ApiTest {
	return New()
}

// Auth Установка Bearer token-а для прохождения авторизации на приватные ручки
func (at *ApiTest) Auth(token string) *ApiTest {
	at.token = token
	return at
}

func (at *ApiTest) Create() *apitest.APITest {
	httpClient := &http.Client{
		Transport: retryTransport{base: http.DefaultTransport},
		Timeout:   30 * time.Second,
	}
	apitestNew := apitest.New().EnableNetworking(httpClient)
	if os.Getenv("DEBUG") == "true" {
		apitestNew = apitestNew.Debug()
	}
	return apitestNew.
		Intercept(func(request *http.Request) {
			request.Header.Set("User-Agent", at.userAgent)

			if err := mergeServiceUrls(at.host, request.URL); err != nil {
				panic(fmt.Sprintf("Не удалось применить хост из API_URL: %s", err.Error()))
			}

			if len(at.token) != 0 {
				request.Header.Add("Authorization", "Bearer "+at.token)
			}
		})
}

// RoundTrip retries only transient HTML errors produced by the reverse proxy.
func (transport retryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var requestBody []byte
	if request.Body != nil {
		var err error
		requestBody, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to buffer request body: %w", err)
		}
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		requestCopy := request.Clone(request.Context())
		if request.Body != nil {
			requestCopy.Body = io.NopCloser(bytes.NewReader(requestBody))
			requestCopy.ContentLength = int64(len(requestBody))
		}

		response, err := transport.base.RoundTrip(requestCopy)
		if err != nil {
			return nil, err
		}
		isTransientHTML := response.StatusCode >= http.StatusInternalServerError &&
			strings.Contains(response.Header.Get("Content-Type"), "text/html")
		if !isTransientHTML || attempt == maxAttempts {
			return response, nil
		}
		response.Body.Close()
		time.Sleep(time.Duration(attempt) * 250 * time.Millisecond)
	}

	return nil, fmt.Errorf("request failed after %d attempts", maxAttempts)
}

func mergeServiceUrls(host string, r *url.URL) error {
	if host == "" {
		return fmt.Errorf("API_URL не задан: проверьте переменные окружения")
	}

	urlParsed, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("host cannot be parsed: %s", err.Error())
	}

	if urlParsed.Path != "" {
		r.Path = path.Join(urlParsed.Path, r.Path)
	}

	r.Scheme = urlParsed.Scheme
	r.Host = urlParsed.Host
	return nil
}
