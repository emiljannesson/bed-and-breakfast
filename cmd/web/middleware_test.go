package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var tH testHandler

	h := NoSurf(&tH)

	switch v := h.(type) {
	case http.Handler:
	// do nothing
	default:
		t.Errorf("type is not http.Handler, type: %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var tH testHandler

	h := SessionLoad(&tH)

	switch v := h.(type) {
	case http.Handler:
	// do nothing
	default:
		t.Errorf("type is not http.Handler, type: %T", v)
	}
}
