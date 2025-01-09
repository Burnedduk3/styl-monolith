package http

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"styl-monolith/pkg/logger"
	"sync"
	"time"
)

// httpClientSingleton is a struct to hold the HTTP client and its configuration
type httpClientSingleton struct {
	client  *http.Client
	baseURL string
}

var (
	httpClientInstance *httpClientSingleton
	once               sync.Once
)

// GetHTTPClient initializes (or retrieves) the singleton HTTP client
func GetHTTPClient(logs *logrus.Logger, baseUrl string) (*http.Client, string) {
	logs.Debug(fmt.Sprintf(logger.CreatingHttpClient, baseUrl))
	once.Do(func() {
		// Initialize the HTTP client
		httpClientInstance = &httpClientSingleton{
			client: &http.Client{
				Timeout: 30 * time.Second, // Set a timeout for requests
			},
			baseURL: baseUrl,
		}
	})

	// Return the configured HTTP client and base URL
	return httpClientInstance.client, baseUrl
}
