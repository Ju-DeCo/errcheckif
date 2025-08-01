package test

import (
	"errors"
	"fmt"
	"os"
)

func mightFail() (string, error) {
	return "hello", errors.New("a sample error")
}

func test() {
	// 错误 1
	_, err := mightFail()
	fmt.Println(err) // 这里仅仅使用，没有检查 err

	// 错误 2 (没有使用 err)
	_, err = mightFail()

	// 正确 1
	_, err = mightFail()
	if err != nil {
		panic(err)
	}

	// 正确 2
	f, err := os.Open("non-existent-file.txt")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")
	}
	defer f.Close()

	// 正确 3
	_, err = mightFail()
	if err == nil {
		// 这种检查虽然不常见，但语法上没错，我们暂时不处理
	}

	// 正确 4
	_, err = mightFail()
	if errors.As(err, &os.ErrNotExist) {
		fmt.Println("file does not exist")
	}

	// 正确 5
	_, _ = mightFail()
}
