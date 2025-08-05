package test

import (
	"errors"
	"fmt"
	"testing"
)

func tfail() error {
	return errors.New("fail")
}

func Test1(t *testing.T) {
	err := tfail()
	fmt.Println(err)
}
