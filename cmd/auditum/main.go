package main

import (
	"os"

	"github.com/infragmo/auditum/internal/cmd/auditum"
)

func main() {
	code := auditum.Execute()
	os.Exit(code)
}
