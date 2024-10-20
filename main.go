// Package main provides a program that periodically sends ATA NOP commands to SATA disks
// to prevent the Load_Cycle_Count from continuously increasing due to IntelliPark.
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	SG_IO = 0x2285
)

type sgIoHdr struct {
	interfaceID    int32
	dxferDirection int32
	cmdLen         uint8
	mxSbLen        uint8
	iovecCount     uint16
	dxferLen       uint32
	dxferp         uintptr
	cmdp           uintptr
	sbp            uintptr
	timeout        uint32
	flags          uint32
	packID         int32
	usrPtr         uintptr
	status         uint8
	maskedStatus   uint8
	msgStatus      uint8
	sbLenWr        uint16
	hostStatus     uint16
	driverStatus   uint16
	resid          int32
	duration       uint32
	info           uint32
}

var verboseFlag *bool

func main() {
	verboseFlag = flag.Bool("verbose", false, "enable verbose output")
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: go-ata-nop <device1> [device2] [device3] ...")
		os.Exit(1)
	}

	devices := flag.Args()

	for {
		for _, device := range devices {
			err := sendNOPCommand(device)
			if err != nil {
				log.Println(err)
			}
		}

		time.Sleep(6 * time.Second)
	}
}

func sendNOPCommand(devicePath string) error {
	if *verboseFlag {
		log.Printf("Sending NOP command to %s\n", devicePath)
	}

	f, err := os.OpenFile(devicePath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %v", devicePath, err)
	}
	defer f.Close()

	cdb := [16]byte{
		0x85,         // ATA PASS-THROUGH (16) command
		(3 << 1) | 1, // Non-data protocol, extend bit set
		0x00,         // ATA NOP command
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0x40, // Set the Command bit in the ATA Control register
		0,
	}

	senseBuffer := [32]byte{}

	hdr := sgIoHdr{
		interfaceID:    'S',
		dxferDirection: -1, // SG_DXFER_NONE
		cmdLen:         uint8(len(cdb)),
		mxSbLen:        uint8(len(senseBuffer)),
		cmdp:           uintptr(unsafe.Pointer(&cdb[0])),
		sbp:            uintptr(unsafe.Pointer(&senseBuffer[0])),
		timeout:        20000,
	}

	// set sgIoHdr values
	buf := make([]byte, unsafe.Sizeof(hdr))
	binary.LittleEndian.PutUint32(buf[0:], uint32(hdr.interfaceID))
	binary.LittleEndian.PutUint32(buf[4:], uint32(hdr.dxferDirection))
	buf[8] = hdr.cmdLen
	buf[9] = hdr.mxSbLen
	binary.LittleEndian.PutUint16(buf[10:], hdr.iovecCount)
	binary.LittleEndian.PutUint32(buf[12:], hdr.dxferLen)
	binary.LittleEndian.PutUint64(buf[16:], uint64(hdr.dxferp))
	binary.LittleEndian.PutUint64(buf[24:], uint64(hdr.cmdp))
	binary.LittleEndian.PutUint64(buf[32:], uint64(hdr.sbp))
	binary.LittleEndian.PutUint32(buf[40:], hdr.timeout)
	binary.LittleEndian.PutUint32(buf[44:], hdr.flags)
	binary.LittleEndian.PutUint32(buf[48:], uint32(hdr.packID))
	binary.LittleEndian.PutUint64(buf[52:], uint64(hdr.usrPtr))

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), SG_IO, uintptr(unsafe.Pointer(&buf[0])))
	if errno != 0 {
		return fmt.Errorf("ioctl failed for %s: %v", devicePath, errno)
	}

	// read results
	status := buf[60]
	if status != 0 {
		return fmt.Errorf("%s: NOP command failed, status = %d", devicePath, status)
	}

	if *verboseFlag {
		log.Printf("NOP command sent to %s\n", devicePath)
	}

	return nil
}
