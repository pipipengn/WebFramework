package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm: Only Support First Level Pointer")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm: Unsupported Expression Type [%v]", expr)
}

func NewErrUnknowField(name string) error {
	return fmt.Errorf("orm: Invalid Field [%s]", name)
}
