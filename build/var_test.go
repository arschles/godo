package build

import (
	"os"
	"testing"

	"github.com/arschles/assert"
	"github.com/arschles/gci/util"
)

func TestVarGetValueDefault(t *testing.T) {
	v := Var{Name: "A", Default: "B"}
	assert.Equal(t, v.Default, v.GetValue(), "value")
}

func TestVarGetValueFromEnv(t *testing.T) {
	name := util.RandIntSuffix("TEST", "_")
	val := "abc"
	assert.NoErr(t, os.Setenv(name, val))
	v := Var{Name: name}
	assert.Equal(t, val, v.GetValue(), "value")
}
