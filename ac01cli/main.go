package main

import (
	"flag"
	"github.com/trumae/ac01"
)

func main() {
	dev := "/dev/ttyUSB0"

	flag.StringVar(&dev, "dev", "/dev/ttyUSB0", "device")
	flag.Parse()

	ac01.DoScan()
}
