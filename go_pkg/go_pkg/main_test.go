package testing_demo

import (
	"testing"
)

func Test(t *testing.T) {
	t.Fatalf("Fatal: %s\n", "FF")
}
