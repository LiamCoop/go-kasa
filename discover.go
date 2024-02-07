package kasa

import (
	"encoding/json"
	"net"
	"strings"
	"time"
)


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
