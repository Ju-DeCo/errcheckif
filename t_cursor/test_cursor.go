package t_cursor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func mightFail() (string, error) {
	return "hello", errors.New("a sample error")
}

func fail() error {
	return errors.New("123")
}

func test() {
	_, err := mightFail()
	fmt.Println(err)

	_, err = mightFail()

	_, err = mightFail()
	if err != nil {
		panic(err)
	}

	f, err := os.Open("non-existent-file.txt")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")
	}
	defer f.Close()

	_, err = mightFail()
	if err == nil {

	}

	_, err = mightFail()
	if errors.As(err, &os.ErrNotExist) {
		fmt.Println("file does not exist")
	}

	_, _ = mightFail()

	if _, err = mightFail(); err != nil {
	}
	if _, err = mightFail(); err == nil {
	}
	if _, err = mightFail(); errors.Is(err, os.ErrNotExist) {
	}
	if _, err = mightFail(); errors.As(err, &os.ErrNotExist) {
	}

	_, err = mightFail()
	if err != nil && err != http.ErrServerClosed {
	}

	_, err = mightFail()
	if err != nil || err != http.ErrServerClosed {
	}

	ctx := context.Background()
	select {
	case <-ctx.Done():
		_, e1 := mightFail()
		if e1 != nil {
		}
	}

	t := 1
	switch t {
	case 1:
		_, e2 := mightFail()
		if e2 != nil {
		}
	}

	if 1 < 2 {
		_, err = mightFail()
	} else {
		_, err = mightFail()
	}
	if err != nil {
	}

	// 并发
	go func() {
		var terr error
		defer func() {
			if terr != nil { //defer中处理
			}
		}()
		terr = fail()
	}()

	var terr error
	go func() {
		terr = fail() // 协程赋值
	}()
	if terr != nil {
	}
}

func error_propagation() (string, error) {
	fail, err := mightFail()
	return fail, err
}

func test_naked_return() (err error) {
	err = errors.New("123")
	return
}