package data

import (
	"testing"
)

func FetchPagesTest(t *testing.T) {
	_, err := FetchPages("", Request{})
	if err != nil {
		t.Error("error is not expected")
	}
}
