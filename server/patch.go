package server

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var dmp = diffmatchpatch.New()

func patch(a, b string) (string, error) {
	patches := dmp.PatchMake(a, b)
	diff := dmp.PatchToText(patches)
	fmt.Println(diff)

	var err error
	patches, err = dmp.PatchFromText(diff)
	if err != nil {
		return "", err
	}

	c, _ := dmp.PatchApply(patches, a)
	return c, nil
}
