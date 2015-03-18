package gohttptest

import (
	"net"
	"net/http"
	"testing"
)

// Server is a simple HTTP server which listens a custom host and port (not on a system-chosen port as in net/http/httptest) for using it in HTTP tests.
type Server struct {
	closer chan struct{}
}

// Stop closes TCP listener on the test HttpTestServer.
func (this *Server) Stop() {
	close(this.closer)
}

// NewServer starts and returns a new Server.
// The caller should call Close when finished, to shut it down.
func NewServer(laddr string, handlers map[string]func(http.ResponseWriter, *http.Request), t *testing.T) *Server {
	srv := &Server{make(chan struct{})}

	go func() {
		mux := http.NewServeMux()
		for k, v := range handlers {
			mux.HandleFunc(k, v)
		}

		server := &http.Server{Handler: mux}
		listener, err := net.Listen("tcp", laddr)
		if err != nil && t != nil {
			t.Error(err.Error())
		}

		select {
		case <- srv.closer:
			listener.Close()
			return
		default:
			//If the channel is still open, continue as normal
		}

		if err := server.Serve(listener); err != nil && t != nil{
			t.Error(err.Error())
		}
	}()
	return srv
}

