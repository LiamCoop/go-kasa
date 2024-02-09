package kasa

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

// returns array of IP addresses, one for each device on the network
func BroadcastAddresses() ([]net.IP, error) {
	var broadcasts []net.IP
	interfaces, err := net.Interfaces()
	if err != nil {
		return broadcasts, err
	}

	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			return broadcasts, err
		}

		for _, addr := range addresses {
			as := addr.String()

			// ignore IPv6 and loopback since kasa devices are v4 only
			if !strings.Contains(as, ":") && !strings.HasPrefix(as, "127.") {
				_, ipnet, err := net.ParseCIDR(as)
				if err != nil {
					return broadcasts, err
				}

				broadcast := net.IP(make([]byte, 4))
				for i := range ipnet.IP {
					broadcast[i] = ipnet.IP[i] | ^ipnet.Mask[i]
				}
				broadcasts = append(broadcasts, broadcast)
			}
		}

	}
	return broadcasts, nil
}

var bufsize = 2048

func Discover(timeout, probes int) (map[string]*Sysinfo, error) {
	m := make(map[string]*Sysinfo)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		return m, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		// XOR the broadcast command
		payload := Scramble(CmdGetSysinfo)
		for i := 0; i < probes; i++ {
			// len(broadcast) yields number of devices on network
			broadcast, _ := BroadcastAddresses()
			for _, b := range broadcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					klogger.Printf("write to udp failed %s", err.Error())
					return
				}
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, bufsize)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			klogger.Printf("Read from UDP failed %s", err.Error())
			break
		}

		res := Unscramble(buffer[:n])

		var kd KasaDevice
		if err = json.Unmarshal(res, &kd); err != nil {
			klogger.Printf("unmarshal: %s\n", err.Error())
			return nil, err
		}

		m[addr.IP.String()] = &kd.GetSysinfo.Sysinfo
	}

	return m, nil
}

func FormatDiscover(m map[string]*Sysinfo) {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Printf("found %d devices\n", len(m))
	for _, k := range keys {
		v := m[k]
		if len(v.Children) == 0 {
			fmt.Printf("%15s: %s %32s [state: %d] [brightness: %3d]\n", k, v.Model, v.Alias, v.RelayState, v.Brightness)
		} else {
			fmt.Printf("%15s: %s %s\n", k, v.Model, v.Alias)
			for _, c := range v.Children {
				fmt.Printf("    ID: %40s%s %26s [state: %d]\n", v.DeviceID, c.ID, c.Alias, c.RelayState)
			}
		}
	}
}
