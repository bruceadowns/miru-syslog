package miru

import (
	"os"
	"testing"
)

func TestOneEvent(t *testing.T) {
	stumptownAddr := os.Getenv("MIRU_STUMPTOWN_ADDR_PORT")
	if stumptownAddr == "" {
		t.Skip("MIRU_STUMPTOWN_ADDR_PORT not found. Skipping test.")
	}

	PostOneEvent(stumptownAddr)
}

func TestManyEvents(t *testing.T) {
	stumptownAddr := os.Getenv("MIRU_STUMPTOWN_ADDR_PORT")
	if stumptownAddr == "" {
		t.Skip("MIRU_STUMPTOWN_ADDR_PORT not found. Skipping test.")
	}

	PostManyEvents(stumptownAddr)
}
