package framework_test

import (
	"testing"

	core "github.com/biuren/hebe/tests/framework/core"
)

func TestMain(m *testing.M) {
	core.RunTestsWithFlags(m)
}
