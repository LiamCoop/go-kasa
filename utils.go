package kasa

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

// This function takes an IP, runs the lookup, then gets sysinfo for that item
func DetailedSystemInfo(ip string) (*KasaDevice, error) {
	d, err := DeviceLookup(ip)
	if err != nil {
		return nil, err
	}
	bytes, err := d.SendTCP(CmdGetSysinfo)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(bytes, &kd); err != nil {
		klogger.Printf("unmarshal: %s", err.Error())
		return nil, err
	}

	return &kd, nil
}

// This function takes an IP, looks up and gets an item back
func DeviceLookup(ip string) (*Device, error) {
	d := Device{IP: ip}
	d.Port = 9999
	d.Parsed = net.ParseIP(ip)

	if d.Parsed == nil {
		// replace with logger usage when this lands in the right place
		fmt.Printf("couldn't parse IP: %s", ip)
		return nil, errors.New("failed to parse provided ip")

	}

	_, err := net.LookupHost(ip)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// use for one-way commands that don't need data back
func (d *Device) SendUDP(cmd string) error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: d.Parsed, Port: d.Port})
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := Scramble(cmd)
	if _, err = conn.Write(payload); err != nil {
		return err
	}

	return nil
}

func (d *Device) SendTCP(cmd string) ([]byte, error) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: d.Parsed, Port: d.Port})
	if err != nil {
		klogger.Printf("Cannot connnect to device: %s", err.Error())
		return nil, err
	}
	defer conn.Close()
	// assume we are on the same LAN, one second is enough
	conn.SetReadDeadline(time.Now().Add(time.Second))

	// send the command with the uint32 "header"
	payload := ScrambleTCP(cmd)
	if _, err = conn.Write(payload); err != nil {
		klogger.Printf("Cannot send command to device: %s", err.Error())
		return nil, err
	}

	// read the uint32 "header" to get the size of the rest of the block
	header := make([]byte, 4)
	n, err := conn.Read(header)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		err := fmt.Errorf("header not 32 bits (4 bytes): %d", n)
		klogger.Printf(err.Error())
		return nil, err
	}
	size := binary.BigEndian.Uint32(header)

	// read the entire rest of the block, then close the connection
	// we could leave the connection open and send subsequent requests
	// but for one-shot, this is enough
	data := make([]byte, size)
	totalread := 0
	for {
		n, err = conn.Read(data[totalread:])
		if err != nil {
			return nil, err
		}
		totalread = totalread + n

		if totalread >= int(size) {
			break
		}
	}

	return Unscramble(data), nil
}
