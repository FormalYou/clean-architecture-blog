package e2e

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/formal-you/clean-architecture-blog/cmd/server/option"
)

var (
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	// Setup: start the server
	router := option.SetupRouter("../../configs")
	testServer = httptest.NewServer(router)
	defer testServer.Close()

	// Setup database and other resources if needed
	// For now, we assume the dev database is running and accessible.

	// Run tests
	code := m.Run()

	// Teardown: clean up resources
	if err := cleanupTestData(); err != nil {
		log.Printf("could not clean up test data: %v", err)
	}

	os.Exit(code)
}

// Placeholder for test data cleanup logic
func cleanupTestData() error {
	// In a real-world scenario, you would connect to the test database
	// and delete the data created during the tests.
	// For example: DELETE FROM users WHERE email LIKE 'testuser%';
	return nil
}

// Helper function to create a request
func newRequest(method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, testServer.URL+url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

// Helper function to create a request with authentication
func newRequestWithAuth(method, url string, body []byte, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, testServer.URL+url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	return client.Do(req)
}
