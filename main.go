package main

import (
	"fmt"

	"github.com/Fawers/battery-info-go/battery"
)

func main() {
	info, err := battery.NewForDefaultDevice()

	if err != nil {
		fmt.Printf("BOOM<%[1]T>: %[1]s\n", err)
	}

	fmt.Printf("%#v\n", info)
}
