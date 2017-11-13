package main

import "flag"

func main() {
	dev := "/dev/ttyUSB0"

	flag.StringVar(&dev, "dev", "/dev/ttyUSB0", "device")
	flag.Parse()

}
