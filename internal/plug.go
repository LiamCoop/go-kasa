package internal

import "fmt"

func (d *Device) SetState(state bool) {
	var command string
	if state {
		command = fmt.Sprintf(CmdSetRelayState, 1)
	} else {
		command = fmt.Sprintf(CmdSetRelayState, 0)
	}

	err := d.SendUDP(command)
	if err != nil {
		Klogger.Printf("Failed to update state")
	}

}
