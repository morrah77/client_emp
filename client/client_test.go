package client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("", "", 0, nil)
	if err != nil {
		t.Fatalf("Could not create client: %s", err.Error())
	}
	if client.maxAttempts != MAX_ATTEMPTS {
		t.Errorf("maxAttempts value is wrong: %v expected, %v actual", MAX_ATTEMPTS, client.maxAttempts)
	}
	var actual, expected uintptr
	expected = reflect.ValueOf(defaultItemValidator).Pointer()
	actual = reflect.ValueOf(client.itemIdValidator).Pointer()
	if actual != expected {
		t.Fatalf("Validator is wrong: %v expected, %v actual", expected, actual)
	}
}

func TestDefaultValidator(t *testing.T) {
	var actual, expected bool
	expected = true
	actual = defaultItemValidator("0 1 2")
	if actual != expected {
		t.Fatalf("Validator works wrong: %v expected, %v actual", expected, actual)
	}
}

func TestClient_FetchList(t *testing.T) {
	expected, testServer, attemptCounter, err := setupTestServer("./test_data/list.json", http.StatusOK, "")
	if err != nil {
		t.Fatalf("Could not set up test server: %s", err.Error())
	}
	defer testServer.Close()

	client, err := NewClient(testServer.URL, testServer.URL, 0, nil)
	if err != nil {
		t.Fatalf("Could not create client: %s", err.Error())
	}
	res, err := client.FetchList()
	if err != nil {
		t.Errorf("Unexpected error got during fetching list: %s", err.Error())
	}
	if bytes.Compare(expected, res) != 0 {
		t.Errorf("Unexpected data got during fetching list:\nexpected: %s\ngot: %s", expected, res)
	}
	if *attemptCounter != 1 {
		t.Errorf("Unexpected attempt count got during fetching list:\nexpected: %v\ngot: %v",
			1, *attemptCounter)
	}
}

func TestClient_FetchItem(t *testing.T) {
	expected, testServer, attemptCounter, err := setupTestServer("./test_data/24.json", http.StatusOK, "")
	if err != nil {
		t.Fatalf("Could not set up test server: %s", err.Error())
	}
	defer testServer.Close()

	client, err := NewClient(testServer.URL, testServer.URL, 0, nil)
	if err != nil {
		t.Fatalf("Could not create client: %s", err.Error())
	}
	res, err := client.FetchList()
	if err != nil {
		t.Errorf("Unexpected error got during fetching item: %s", err.Error())
	}
	if bytes.Compare(expected, res) != 0 {
		t.Errorf("Unexpected data got during fetching item:\nexpected: %s\ngot: %s", expected, res)
	}
	if *attemptCounter != 1 {
		t.Errorf("Unexpected attempt count got during fetching item:\nexpected: %v\ngot: %v",
			1, *attemptCounter)
	}
}

func TestClient_FetchItemWithBadRequestStatus(t *testing.T) {
	expected, testServer, attemptCounter, err := setupTestServer("./test_data/400.json", http.StatusBadRequest, "")
	if err != nil {
		t.Fatalf("Could not set up test server: %s", err.Error())
	}
	defer testServer.Close()

	client, err := NewClient(testServer.URL, testServer.URL, 0, nil)
	if err != nil {
		t.Fatalf("Could not create client: %s", err.Error())
	}
	res, err := client.FetchList()
	expectedErrorMessage := "400 Bad Request"
	if (err == nil) || (err.Error() != expectedErrorMessage) {
		t.Errorf("Unexpected error got during fetching item with 400:\nexpected error with message: %s\ngot: %s", expectedErrorMessage, err)
	}
	if bytes.Compare(expected, res) != 0 {
		t.Errorf("Unexpected data got during fetching item with 400:\nexpected: %s\ngot: %s", expected, res)
	}
	if *attemptCounter != client.maxAttempts {
		t.Errorf("Unexpected attempt count got during fetching item with 400:\nexpected: %v\ngot: %v",
			client.maxAttempts, *attemptCounter)
	}
}

// TODO add test case for other server errors
func TestClient_FetchItemWithError(t *testing.T) {
	_, testServer, _, err := setupTestServer("./test_data/400.json", http.StatusPermanentRedirect, "Location:/")
	if err != nil {
		t.Fatalf("Could not set up test server: %s", err.Error())
	}
	defer testServer.Close()

	client, err := NewClient(testServer.URL, testServer.URL, 0, nil)
	if err != nil {
		t.Fatalf("Could not create client: %s", err.Error())
	}
	expectedError := "Get /: stopped after 10 redirects"
	_, err = client.FetchList()
	if (err == nil) || (err.Error() != expectedError) {
		t.Errorf("Expected error didn't get during fetching item with error:\nexpected error with message %s\ngot %#v",
			expectedError, err)
	}
}

// TODO consider returning structure instead of result set
func setupTestServer(dataFilePath string, statusToRespond int, additionalHeader string) (expected []byte,
	server *httptest.Server,
	attemptCounter *int, err error) {
	b, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		return nil, nil, nil, errors.New(fmt.Sprintf("Could not read from file: %s", err.Error()))
	}
	expected = make([]byte, len(b))
	copiedBytes := copy(expected, b)
	if copiedBytes != len(b) {
		return nil, nil, nil, errors.New(fmt.Sprintf("Inappropriate length!\nexpected: %v\ngot: %v", len(b),
			copiedBytes))
	}
	ac := 0
	attemptCounter = &ac

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*attemptCounter++
		if len(additionalHeader) > 0 {
			parts := strings.Split(additionalHeader, ":")
			if len(parts) < 2 {
				w.Header().Add(parts[0], "")
			}
			w.Header().Add(parts[0], parts[1])
		}
		w.WriteHeader(statusToRespond)
		fmt.Fprint(w, string(b))
	}))
	return expected, testServer, attemptCounter, err
}
