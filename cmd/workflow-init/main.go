package main

import (
	"fmt"

	"github.com/samdoiron/workflow"
)

func main() {
	workflow.ResetTables()
	fmt.Println("Tables have been reset")
}
