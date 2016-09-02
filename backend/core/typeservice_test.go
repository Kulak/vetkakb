package core

import (
	"log"
	"testing"
)

func TestHtmlToPlain(t *testing.T) {
	prov := htmlProvider()
	html := "<p>test line</p><ol><li>one</li><li>two</li></ol>"
	plain, err := prov.ToPlain([]byte(html))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	log.Print(plain)
}
