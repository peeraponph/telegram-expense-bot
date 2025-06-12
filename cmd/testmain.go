package main

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func testmain() {
    client := gosseract.NewClient()
    defer client.Close()
    fmt.Println("🧠 NewClient() is accessible — version:", client.Version())
}
