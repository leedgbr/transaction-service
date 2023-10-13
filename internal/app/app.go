package app

import (
	"context"
	"fmt"
	"net/http"
)

// New sets up a new application (to manage the http server) with the supplied Dependencies
func New(deps Dependencies) *App {
	router := newRouter(deps)
	return &App{
		Router: router,
		Server: newServer(router),
	}
}

// App represents the 'application'
type App struct {
	Router http.Handler
	Server *http.Server
}

// Start launches the http web application on the provided port
func (a *App) Start(port int) error {
	fmt.Println("Using port:", port)
	listener, err := NewListener(port)
	if err != nil {
		return err
	}
	return http.Serve(listener, a.Router)
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) {
	a.Server.Shutdown(ctx)
}
