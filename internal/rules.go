package internal

import "encoding/json"

func (d *Device) GetRules() (Countdown, error) {
	command := CmdGetCountdownRules
	bytes, err := d.SendTCP(command)
	if err != nil {
		Klogger.Printf("failed to send TCP")
	}

	var kd KasaDevice
	if err = json.Unmarshal(bytes, &kd); err != nil {
		Klogger.Printf("unmarshal: %s", err.Error())
		return kd.Countdown, err
	}

	return kd.Countdown, nil
}
