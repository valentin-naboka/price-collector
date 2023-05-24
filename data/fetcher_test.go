package data

import (
	"testing"
)

func TestFetchPages(t *testing.T) {
	_, err := FetchPages("", Request{})
	if err == nil {
		t.Error("error is expected")
	}
}
