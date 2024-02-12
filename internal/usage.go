package internal

import (
	"encoding/json"
	"fmt"
)

// CmdGetMonthlyUsageInYear
func (d *Device) Usage(usageYear int16) (Schedule, error) {
	command := fmt.Sprintf(CmdGetMonthlyUsageInYear, usageYear)
	bytes, err := d.SendTCP(command)

	if err != nil {
		Klogger.Printf("failed to send TCP")
	}

	var kd KasaDevice
	if err = json.Unmarshal(bytes, &kd); err != nil {
		Klogger.Printf("unmarshal: %s", err.Error())
		return kd.Schedule, err
	}

	return kd.Schedule, nil
}

// Does not work
// CmdLEDOff
func (d *Device) DeviceOff() error {
	command := fmt.Sprintf(CmdLEDOff, 0)
	fmt.Println(command)
	err := d.SendUDP(command)
	if err != nil {
		Klogger.Printf("failed to send UDP")
	}
	return nil
}
