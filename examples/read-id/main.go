package main

import (
	"fmt"

	"github.com/holoplot/go-swd/pkg/io/bitbang"
	"github.com/holoplot/go-swd/pkg/swd"
)

func main() {
	linuxGPIO, err := bitbang.NewLinuxGPIO("/dev/gpiochip0", 17, 27, 100000)
	if err != nil {
		panic(err)
	}

	swd := swd.New(linuxGPIO)

	idcode, err := swd.IDCode()
	if err != nil {
		panic(err)
	}

	fmt.Printf("IDCODE: 0x%08x\n", idcode)
}
