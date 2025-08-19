package test

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func tfail() error {
	return errors.New("fail")
}

func Test1(t *testing.T) {
	err := tfail()
	fmt.Println(err)
}

func Test01(t *testing.T) {
	var err error

	if 1 < 2 {
		err = tfail()
	} else {
		_, err = os.Open("test.txt")
	}
	fmt.Println(err)
}
