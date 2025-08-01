package main

import (
	"github.com/Ju-DeCo/errcheckif"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(errcheckif.Analyzer)
}
