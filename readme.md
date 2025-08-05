# errcheckif

如果函数调用返回值包含`error`类型，那么这个`error`变量 err 必须在后续 `if` 语句中被检查，检查条件可以是：

* err != nil
* err == nil
* errors.Is(err, ***)
* errors.As(err, ***)

或者通过 `return` 进行错误传递。

默认跳过测试文件（以`_test.go`结尾）。

### 例子：

``` go
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

func error_propagation() (string, error) {
    // 正确 5 错误传递
    fail, err := mightFail()
    return fail, err
}

// 正确6 裸返回
func test_naked_return() (err error) {
    err = errors.New("123")
    return
}
```

### 进行测试：
```
go mod tidy
go run .\cmd\errcheckif\ .\test\test.go
```

### 局限性

**控制流误报**
``` go
if 1 < 2 {
    _, err = mightFail()
} else {
    _, err = mightFail()
}
if err != nil {
}

```

**并发误报**
``` go
go func() {
    var terr error
    defer func() {
        if terr != nil {
        }
    }()
    terr = fail()
}()
```
