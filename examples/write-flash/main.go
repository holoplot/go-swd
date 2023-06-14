package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/holoplot/go-swd/pkg/io/bitbang"
	"github.com/holoplot/go-swd/pkg/stm32"
	"github.com/holoplot/go-swd/pkg/swd"
)

func main() {
	fileFlag := flag.String("file", "", "file to read content from")
	flag.Parse()

	if *fileFlag == "" {
		panic("no file specified")
	}

	content, err := ioutil.ReadFile(*fileFlag)
	if err != nil {
		panic(err)
	}

	linuxGPIO, err := bitbang.NewLinuxGPIO("/dev/gpiochip1", 81, 80, 1000000)
	if err != nil {
		panic(err)
	}

	if err := linuxGPIO.LineReset(); err != nil {
		panic(err)
	}

	s := swd.New(linuxGPIO)

	if err := s.PowerOnReset(); err != nil {
		panic(err)
	}

	if err := s.Abort(swd.AbortAllFlags()); err != nil {
		panic(err)
	}

	if err := s.Select(0, 0); err != nil {
		panic(err)
	}

	stm := stm32.New(s)
	flash := stm.Flash()

	log.Printf("Erasing flash...")

	if err := flash.EraseAll(time.Minute); err != nil {
		panic(err)
	}

	stat, err := os.Stat(*fileFlag)
	if err != nil {
		panic(err)
	}

	log.Printf("Writing flash (%d bytes) ...", stat.Size())

	buf := bytes.NewBuffer(content)

	if err := flash.Write(0, buf); err != nil {
		panic(err)
	}

	log.Printf("Verifying flash...")

	buf = bytes.NewBuffer(nil)

	if err := flash.Read(0, uint32(len(content)), buf); err != nil {
		panic(err)
	}

	if !bytes.Equal(buf.Bytes(), content) {
		panic("verification failed")
	}

	log.Printf("Resetting...")

	if err := stm.RunAfterReset(); err != nil {
		panic(err)
	}

	if err := stm.Reset(); err != nil {
		panic(err)
	}
}
