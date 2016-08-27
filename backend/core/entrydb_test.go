package core_test

import (
	"testing"

	"horse.lan.gnezdovi.com/vetkakb/backend/core"
)

func TestInitJournalTest(t *testing.T) {
	_, err := core.LoadConfig("test.gcfg")
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

}
