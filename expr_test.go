package expr_test

import (
	"bytes"
	"testing"

	"github.com/grailbio/testutil/expect"
	expr "github.com/yasushi-saito/go-expr"
)

func eval(t *testing.T, s string) float64 {
	in := bytes.NewReader([]byte(s))
	inst := expr.Compile(in)
	return expr.Run(inst)
}

func TestBasic(t *testing.T) {
	expect.EQ(t, eval(t, "1+2"), 3.0)
}
