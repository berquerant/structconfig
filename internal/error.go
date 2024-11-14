package internal

import (
	"errors"
	"fmt"
)

var (
	ErrStructConfig     = errors.New("StructConfig")
	ErrNotStruct        = errors.New("NotStruct")
	ErrNotStructPointer = errors.New("NotStructPointer")
)

func Errorf(format string, v ...any) error {
	msg := fmt.Sprintf(format, v...)
	return fmt.Errorf("%w: %s", ErrStructConfig, msg)
}

func JoinErrors(err ...error) error {
	return errors.Join(append([]error{ErrStructConfig}, err...)...)
}
