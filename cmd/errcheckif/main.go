// cmd/errcheckif/main.go
package main

import (
	"github.com/Ju-DeCo/errcheckif" // 替换为你的模块路径
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(errcheckif.Analyzer)
}
