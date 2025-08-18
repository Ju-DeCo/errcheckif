package test

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
	// 错误 1
	_, err := mightFail()
	fmt.Println(err) // 这里仅仅使用，没有检查 err

	// 错误 2 (没有使用 err)
	_, err = mightFail()

	// 错误 3 直接忽略错误
	_, _ = mightFail()
	_ = fail()

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

	// 正确 5 if-init模式
	if _, err = mightFail(); err != nil {
	}
	if _, err = mightFail(); err == nil {
	}
	if _, err = mightFail(); errors.Is(err, os.ErrNotExist) {
	}
	if _, err = mightFail(); errors.As(err, &os.ErrNotExist) {
	}

	// 正确 6 逻辑与 与 逻辑或
	_, err = mightFail()
	if err != nil && err != http.ErrServerClosed {
	}

	_, err = mightFail()
	if err != nil || err != http.ErrServerClosed {
	}

	// 正确 7 select 与 switch 语句
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

	// 未能解决的问题

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
	// 正确 8 错误传递
	fail, err := mightFail()
	return fail, err
}

// 正确 9 裸返回 naked return
func test_naked_return() (err error) {
	err = errors.New("123")
	return
}

// 错误
func test_cross(cond bool) {
	err := fail() // Linter 发现 err A
	if cond {
		err = fail() // 一个新的 err B
		if err != nil {
			return
		}
	}
	fmt.Println(err)
	// 此处 err A 未被处理
}

// ================= ifelse  ==================
func rterr() error {
	return errors.New("test")
}

// 正确 1
func ttest01(cond bool) {
	var err error

	if cond {
		err = rterr()
	} else {
		_, err = os.Open("test.txt")
	}
	if err != nil {
	}
}

// 正确 2
func ttest02(cond bool) {
	var err error

	if cond {
		err = rterr()
	} else {
		_, err = os.Open("test.txt")
	}
	if err == nil {
	}
}

// 正确 3
func ttest03(cond bool) {
	var err error

	if cond {
		err = rterr()
	} else {
		_, err = os.Open("test.txt")
	}
	if errors.Is(err, os.ErrNotExist) {
	}
}

// 正确 4
func ttest04(cond bool) {
	var err error

	if cond {
		err = rterr()
	} else {
		_, err = os.Open("test.txt")
	}
	if errors.As(err, os.ErrNotExist) {
	}
}

// 错误 1
func ftest01(cond bool) {
	var err error

	if cond {
		err = rterr()
	} else {
		_, err = os.Open("test.txt")
	}
	fmt.Println(err)
}

// 错误 2
func ftest02(cond bool) {
	err := rterr()
	if err != nil {
	}

	if cond {
		err = rterr()
	} else {
		_, err = os.Open("test.txt")
	}
}
