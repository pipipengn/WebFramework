package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm: Only Support First Level Pointer")
	ErrNoRows      = errors.New("orm: No Data")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm: Unsupported Expression Type [%v]", expr)
}

func NewErrUnknowField(name string) error {
	return fmt.Errorf("orm: Unknown Field [%s]", name)
}

func NewErrUnknowColumn(col string) error {
	return fmt.Errorf("orm: Unknown Column [%s]", col)
}

func NewErrInvalidTagContent(key string) error {
	return fmt.Errorf("orm: Invalid Tag Content [%s]", key)
}
