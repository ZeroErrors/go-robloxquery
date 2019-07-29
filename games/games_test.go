package list

import (
	"testing"
)

func TestDo(t *testing.T) {
	// TODO: Build out some testing
	resp, err := Do(Request{UniverseIDs: []int64{383310974, 88070565}})
	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}
