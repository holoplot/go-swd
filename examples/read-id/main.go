package main

import (
	"fmt"

	"github.com/holoplot/go-swd/pkg/io/bitbang"
	"github.com/holoplot/go-swd/pkg/swd"
)

func main() {
	linuxGPIO, err := bitbang.NewLinuxGPIO("/dev/gpiochip1", 81, 80, 1000000)
	if err != nil {
		panic(err)
	}

	if err := linuxGPIO.LineReset(); err != nil {
		panic(err)
	}

	swd := swd.New(linuxGPIO)

	idcode, err := swd.IDCode()
	if err != nil {
		fmt.Println(err)
	}

	idcode, err = swd.IDCode()
	if err != nil {
		panic(err)
	}

	fmt.Printf("IDCODE: 0x%08x\n", idcode)
}
