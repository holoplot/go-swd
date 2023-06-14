# Pure Go library for the Serial Wire Debug (SWD) protocol

Serial Wire Debug (SWD) is a 2-pin (SWDIO/SWCLK) electrical alternative JTAG interface that has the
JTAG protocol on top. It was developed by ARM® and is widely used for programming/debugging
microcontrollers, PMICs, and other embedded devices. More information on SWD can be found in the
[ARM Debug Interface v5 Architecture Specification](https://developer.arm.com/documentation/ihi0031/d/).

This library aims to implement the SWD protocol in pure Go. It was started to program the flash
of STM32 MCUs.

# Structure

The library is split into distinct parts that can be used independently.

## Low-level IO

The `io` package provides the low-level SWD interface that can be implemented by different
transport implementations. The library provides bitbang implementations for generic Linux
sysfs GPIOs as well as Raspberry Pi GPIOs.

Hardware accelerated transports may implement the `io.Accessor` interface directly.

## SWD protocol

The `swd` package contains the SWD protocol implementation for accessing debug and access
port registers such as CSW, TAR or DRW.

## Core Debug

The Core Debug layer is a higher-level interface that allows to access the debug registers
or memory of a Cortex-M MCU. It is implemented in the `core-debug` package.

For more information on the Core Debug interface, refer to the
[Cortex-M3 Technical Reference Manual r1p1](https://developer.arm.com/documentation/ddi0337/e/).

## STM32

This layer provides convenience functions for interacting with STM32 MCUs such as reading,
writing and erasing flash memory. It is implemented in the `stm32` package.

# Examples

Please refer to the `examples` directory for simple examples that read the IDCODE of a
STM32 MCU or program the flash memory.

# Contributing

Contributions are welcome as this project was started with a small scope but aims for
more functionality and device support. Please open pull requests to extend it.

# License

This library is licensed under the MIT license. See the LICENSE file for more information.
