

Understanding TP-Link's Kasa Discover

Firstly UDP broadcast packet is sent out to 255.255.255.255:9999
The application sends out the following content, 
encrypted with an XOR encryption

| {"system":{"get_sysinfo":{}}}
