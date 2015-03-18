package gohttptest

import (
	"net/http"
	"testing"
	"net/url"
	"io/ioutil"
)

const (
	endpoint = "http://localhost:8081/"
	api_v1 = 1
	api_v2 = 2
)

// Endpoint
type api struct {
	path    string
	payload string
}

var apis = map[int]api{
	api_v1: api{"/api/v1", `{"message":"api.v1"}`},
	api_v2: api{"/api/v2", `{"message":"api.v2"}`},
}

func TestServer(t *testing.T) {
	srv := setUp(t)
	defer tearDown(srv)
	
	// Check apis
	for k, _ := range apis {
		checkApi(t, k)
	}
}

func setUp(t *testing.T) *Server {
	handlers := map[string]func (http.ResponseWriter, *http.Request){
		apis[1].path : HandlerApi(apis[1].payload),
		apis[2].path : HandlerApi(apis[2].payload),
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

func HandlerApi(payload string) (func(http.ResponseWriter, *http.Request)) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-type", "application/json")
		writer.Write([]byte(payload))
	}
}

func checkApi(t *testing.T, api int) {
	var (
		req  *http.Request
		err  error
		resp *http.Response
	)
	
	if req, err = http.NewRequest("GET", endpoint + apis[api].path, nil); err != nil {
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
	if apis[api].payload != bodyString {
		t.Errorf("Responce isn't correct: expected %s, got %s", apis[api].payload, bodyString)
	}
}
