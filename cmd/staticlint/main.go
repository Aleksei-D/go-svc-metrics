package main

import (
	"go-svc-metrics/internal/linters"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(linters.NewCheckers()...)
}
