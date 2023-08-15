package main

import (
	"github.com/emiljannesson/bed-and-breakfast/internal/config"
	"github.com/go-chi/chi"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	// do nothing, test passed
	default:
		t.Errorf("type is not *chi.Mux, type: %T", v)
	}
}
