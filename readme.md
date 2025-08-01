# errcheckif

用于检测函数调用返回了err，但是没有进行检测的情况，包含：
* `err != nil`
* `err == nil`
* `errors.Is`
* `errors.As`

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

// 正确 5
_, _ = mightFail()

// 正确 6 if-init模式
if _, err = mightFail(); err != nil {
}
if _, err = mightFail(); err == nil {
}
if _, err = mightFail(); errors.Is(err, os.ErrNotExist) {
}
if _, err = mightFail(); errors.As(err, &os.ErrNotExist) {
}
```

### 进行测试：
```
go mod tidy
go run .\cmd\errcheckif\ .\test\test.go
```