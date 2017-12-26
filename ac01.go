package ac01

import (
	"fmt"
	"log"
	"time"

	"encoding/hex"

	"github.com/snksoft/crc"
	"github.com/tarm/serial"
)

func NewMsgGetReaderInformation() []byte {
	ret := []byte{}

	ret = append(ret, 0xBB) //preamble
	ret = append(ret, 0x00) //msg type
	ret = append(ret, 0x03) //code
	ret = append(ret, 0x00) //PL(MSB)
	ret = append(ret, 0x01) //PL(LSB)
	ret = append(ret, 0x02) //Arg
	ret = append(ret, 0x7E) //endmark

	tmp := make([]byte, len(ret))
	copy(tmp, ret)

	rcrc := crc.CalculateCRC(crc.CCITT, tmp[1:])

	var h, l uint8 = uint8(rcrc >> 8), uint8(rcrc & 0xff)
	ret = append(ret, h)
	ret = append(ret, l)
	return ret
}

func NewMsgStartRead() []byte {
	ret := []byte{}

	ret = append(ret, 0xBB) //preamble
	ret = append(ret, 0x00) //msg type
	ret = append(ret, 0x36) //code
	ret = append(ret, 0x00) //PL(MSB)
	ret = append(ret, 0x05) //PL(LSB)
	ret = append(ret, 0x02) //Reserve
	ret = append(ret, 0x00) //MTNU
	ret = append(ret, 0x00) //MTIME in seconds
	ret = append(ret, 0x00) //RC(MSB)
	ret = append(ret, 0x64) //RC(LSB)
	ret = append(ret, 0x7E) //endmark

	tmp := make([]byte, len(ret))
	copy(tmp, ret)

	rcrc := crc.CalculateCRC(crc.CCITT, tmp[1:])
	var h, l uint8 = uint8(rcrc >> 8), uint8(rcrc & 0xff)
	ret = append(ret, h)
	ret = append(ret, l)
	return ret
}

func NewMsgStartRead2() []byte {
	ret := []byte{
		0xBB, 0x00, 0x36, 0x00, 0x05,
		0x02, 0x00, 0x00, 0x00, 0x64,
		0x7E}

	tmp := make([]byte, len(ret))
	copy(tmp, ret)

	rcrc := crc.CalculateCRC(crc.CCITT, tmp[1:])
	var h, l uint8 = uint8(rcrc >> 8), uint8(rcrc & 0xff)
	ret = append(ret, h)
	ret = append(ret, l)
	return ret
}

func NewMsgStopRead() []byte {
	ret := []byte{}

	ret = append(ret, 0xBB) //preamble
	ret = append(ret, 0x00) //msg type
	ret = append(ret, 0x37) //code
	ret = append(ret, 0x00) //PL(MSB)
	ret = append(ret, 0x00) //PL(LSB)
	ret = append(ret, 0x7E) //endmark

	tmp := make([]byte, len(ret))
	copy(tmp, ret)

	rcrc := crc.CalculateCRC(crc.CCITT, tmp[1:])
	var h, l uint8 = uint8(rcrc >> 8), uint8(rcrc & 0xff)
	ret = append(ret, h)
	ret = append(ret, l)
	return ret
}

func SendGetReaderInformation() {
	msg := NewMsgGetReaderInformation()

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", buf[:n])
}

func printHex(bs []byte) {
	fmt.Println(hex.Dump(bs))
}

func DoScan(sec int) {
	mstart := NewMsgStartRead()
	mstop := NewMsgStopRead()

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	read := true

	go func() {
		timer := time.NewTimer(time.Second * time.Duration(sec))
		<-timer.C
		log.Println("Timer expired")

		printHex(mstop)
		_, err := s.Write(mstop)
		if err != nil {
			log.Fatal(err)
		}
		read = false
	}()

	printHex(mstart)
	n, err := s.Write(mstart)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if !read {
			break
		}

		buf := make([]byte, 128)
		n, err = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		printHex(buf[:n])
	}
}
