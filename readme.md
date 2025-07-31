# errcheckif

用于检测函数调用返回了err，但是没有进行检测的情况，包含三种情况：
* err != nil
* errors.Is
* errors.As

例子：

``` go
// 错误 1
_, err := mightFail()
fmt.Println(err) // <- 这里仅仅使用，没有检查 err

// 错误 2 (没有使用 err)
_, err = mightFail()

// 正确的例子 1
_, err = mightFail()
if err != nil {
    panic(err)
}

// 正确的例子 2
f, err := os.Open("non-existent-file.txt")
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("file does not exist")
}
defer f.Close()

// 正确的例子 3
_, err = mightFail()
if err == nil {
    // 这种检查虽然不常见，但语法上没错，我们暂时不处理
}

// 正确的例子 4
_, err = mightFail()
if errors.As(err, &os.ErrNotExist) {
    fmt.Println("file does not exist")
}

```