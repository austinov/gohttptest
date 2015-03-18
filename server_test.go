package gohttptest

import (
	"net/http"
	"testing"
	"net/url"
	"io/ioutil"
	"fmt"
)

const (
	endpoint = "http://localhost:8081/tests/events"
	payload = `{"message": {"id": 1, "name": "Mike"}}`
)

func TestServer(t *testing.T) {
	srv := setUp(t)
	defer tearDown(srv)
	checkApi(t)
}

func setUp(t *testing.T) *Server {
	handlers := map[string]func (http.ResponseWriter, *http.Request){
		"/" : getHandlers(t),
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		t.Error(err)
	}

	return NewServer(u.Host, handlers, t)
}

func tearDown(srv *Server) {
	if srv != nil {
		srv.Stop()
	}
}

// Setup handlers for requests from a client
func getHandlers(t *testing.T) (func(http.ResponseWriter, *http.Request)) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			http.Error(writer, "Bad request", http.StatusBadRequest)
			return
		}
		if request.URL.Path == "/tests/events" {
			getEvents(t, writer, request)
		} else {
			http.NotFound(writer, request)
		}
	}
}

// Helper emulates the API method on the server
func getEvents(t *testing.T, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")
	writer.Write([]byte(payload))
}

func checkApi(t *testing.T) {
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	
	if req, err = http.NewRequest("GET", endpoint, nil); err != nil {
		t.Error(err)
		return
	}
    client := &http.Client{}

	if resp, err = client.Do(req); err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	fmt.Printf("Result: %s\n", bodyString)
	if payload != bodyString {
		t.Errorf("Responce isn't correct: expected %s, got %s", payload, bodyString)
	}
}
