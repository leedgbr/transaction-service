package fixture

import (
	"context"
	"net"
	"net/http/httptest"
	"testing"

	"transaction-service/internal/app"
	"transaction-service/internal/id"
)

const nextAvailablePort = 0

// NewTestServer creates a test http server that is convenient for running the application for the purpose of
// integration testing.
func NewTestServer() *TestServer {
	return &TestServer{}
}

// TestServer exposes the BaseURL on which the TestServer is made available, and contains the 'application'.
type TestServer struct {
	BaseURL     string
	application *app.App
}

// Start wires up the test server with the minimum wiring modified with specifically stubbed dependencies to make
// integration testing easier.  This includes a transaction.IDGenerator that produces predictable IDs so that we are
// able to reference them in test scenarios.  It also includes a stub http client which will deliver configured stub
// responses to http calls to apis on which this transaction-service depends.  A port number of '0' is provided which
// results in the next available port being allocated to the test server.  There should never be port conflicts with
// any running integration tests or standalone server.
func (s *TestServer) Start(t *testing.T) {
	s.application = app.New(app.NewDependencies(id.NewSequentialGenerator(), NewStubHttpClient()))
	server := httptest.NewUnstartedServer(s.application.Router)
	listener, err := app.NewListener(nextAvailablePort)
	if err != nil {
		t.Fatal(err)
	}
	server.Listener = listener
	server.Start()
	s.BaseURL = BaseURL(Port(listener))
}

// Stop stops the test server.
func (s *TestServer) Stop() {
	s.application.Stop(context.Background())
}

// Port returns the current port assigned to the provided net.Listener.
func Port(listener net.Listener) int {
	return listener.Addr().(*net.TCPAddr).Port
}
