package kasa

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func Usage(timeout, probes int) (map[string]*Sysinfo, error) {

    conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
    if err != nil {
        
        return nil, err
    }
    defer conn.Close()

    conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

    go func() {
        // XOR the broadcast command
        payload := Scramble(CmdGetMonthlyUsageInYear)
        for i:= 0; i < probes; i++ {
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
        n, _, err := conn.ReadFromUDP(buffer)
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

        str, err := json.Marshal(kd.Schedule)
        if err != nil {
            klogger.Printf("marshal failed on month list: %d")
            return nil, err
        }

        fmt.Printf("marshal %s", str)
    }

    return nil, nil
}
