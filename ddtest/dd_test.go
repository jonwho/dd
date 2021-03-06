package ddtest

import (
	"testing"
)

func TestWSJ(t *testing.T) {
	content, _ := WSJ("aapl")

	if len(content) < 1 {
		t.Error("Expected at least some articles to fill content slice since not all articles are behind a paywall")
	}
}

func TestReadyHoset(t *testing.T) {
	host := ReadHost("")
	if host != "" {
		t.Errorf("Expected ...?\nActual: %s", host)
	}
}
