package main

import (
	"fmt"

	"github.com/Fawers/battery-info-go/battery"
)

func main() {
	info, err := battery.NewFromDefaultDevice()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", info)
}
