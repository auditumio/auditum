package main

import (
	"os"

	"github.com/auditumio/auditum/internal/cmd/auditum"
)

func main() {
	code := auditum.Execute()
	os.Exit(code)
}
