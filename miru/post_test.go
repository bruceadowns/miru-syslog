package miru

import (
	"os"
	"testing"
)

func TestOneEvent(t *testing.T) {
	stumptownAddr := os.Getenv("MIRU_STUMPTOWN_HOST_PORT")
	if len(stumptownAddr) == 0 {
		t.Skip("MIRU_STUMPTOWN_HOST_PORT not found. Skipping test.")
	}

	PostOneEvent(stumptownAddr)
}

func TestManyEvents(t *testing.T) {
	stumptownAddr := os.Getenv("MIRU_STUMPTOWN_HOST_PORT")
	if len(stumptownAddr) == 0 {
		t.Skip("MIRU_STUMPTOWN_HOST_PORT not found. Skipping test.")
	}

	PostManyEvents(stumptownAddr)
}
