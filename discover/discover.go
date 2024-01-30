package discover

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
    "github.com/liamcoop/go-kasa/xor"
)



var DISCOVERY_QUERY = map[string]map[string]interface{}{
    "system": {"get_sysinfo": nil},
}
var DISCOVERY_PORT = 9999

var DISCOVERY_QUERY_2_ENCODED = "020000010000000000000000463cb5d3"
var DISCOVERY_PORT_2 = 20002

var DISCOVERY_PACKETS = 3

func Discover() {
    target := "255.255.255.255"
    // discover_packets := 3
    // discover_timeout := 5

    serverAddr, err := net.ResolveUDPAddr("udp", target + ":" + strconv.Itoa(DISCOVERY_PORT))
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        return
    }

    conn, err := net.DialUDP("udp", nil, serverAddr)

    if err != nil {
        fmt.Println("Error connecting to UDP server:", err)
        return
    }
    defer conn.Close()

    
    req, err := json.Marshal(DISCOVERY_QUERY)

    if err != nil {
        fmt.Println("marshalling failed")
    }

    // need to figure out xor encryption
    encrypted_req := xor.EncryptDecrypt(req, "")

    for i := 1; i <= DISCOVERY_PACKETS; i++ {
        conn.Write(encrypted_req[4:], )
    }

    // what to send to server to get it to give up the goods
    // message := []byte("Hello, UDP Server!")

    // _, err = conn.Write(message)
    /*
    if err != nil {
        fmt.Println("Error sending message:", err)
        return
    }

    fmt.Println("Message sent to UDP server", string(message))
    */

    
}
