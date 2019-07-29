package list

import (
	"testing"
)

func TestDo(t *testing.T) {
	// TODO: Build out some testing
	resp, err := Do(Request{})
	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}
