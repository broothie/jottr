package server

import "testing"

func Test_patch(t *testing.T) {
	c, err := patch("andrew booth", "andy both")
	if err != nil {
		t.Error(err)
	}

	t.Error(c)
}
