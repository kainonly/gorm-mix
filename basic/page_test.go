package basic

import "testing"

func TestGeneratePage(t *testing.T) {
	if err := GeneratePage(db); err != nil {
		t.Error(err)
	}
}