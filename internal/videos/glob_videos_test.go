package videos

import (
	"fmt"
	"testing"
)

func TestJavWalk(t *testing.T) {
	err := JavWalk("F:Download/JAV/JAV_output")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(FilesList)
}
