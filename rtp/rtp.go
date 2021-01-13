package rtp

import (
	"fmt"
	"net"
	"strconv"
)

/*
	    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |V=2|P|X|  CC   |M|     PT      |       sequence number         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           timestamp                           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           synchronization source (SSRC) identifier            |
   +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
   |            contributing source (CSRC) identifiers             |
   |                             ....                              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

// Listen - listening server
func Listen(port int) {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error binding port!")
	}

	ServerConn, _ := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()

	buf := make([]byte, 4096)
	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error packet")
			continue
		}

		cpy := make([]byte, n)
		copy(cpy, buf)
		go handleBuf(cpy)
	}
}

func handleBuf(buf []byte) {
	Version := buf[0] & 0x03
	// Padding := buf[0]&1<<2 != 0
	// Ext := buf[0]&1<<3 != 0
	// CSRC := make([]uint, buf[0]>>4)
	// Marker := buf[1]&1 != 0
	// PayloadType := buf[1] >> 1
	// SequenceNumber := toUint(buf[2:4])
	// Timestamp := toUint(buf[4:8])
	// SyncSource := toUint(buf[8:12])
	// fmt.Println(Version, Padding, Ext, CSRC, Marker, PayloadType, SequenceNumber, Timestamp, SyncSource)
	fmt.Println(Version)
}

func toUint(arr []byte) (ret uint) {
	for i, b := range arr {
		ret |= uint(b) << (8 * uint(len(arr)-i-1))
	}
	return ret
}
