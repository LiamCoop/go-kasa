package internal

import (
	"encoding/json"
	"fmt"
)

func (d *Device) SetAlias(a string) error {
	command := fmt.Sprintf(CmdDeviceAlias, a)
	bytes, err := d.SendTCP(command)
	if err != nil {
		Klogger.Printf("failed to send TCP")
	}

	var kd KasaDevice
	if err = json.Unmarshal(bytes, &kd); err != nil {
		Klogger.Printf("unmarshal: %s", err.Error())
		return err
	}

	return nil
}
