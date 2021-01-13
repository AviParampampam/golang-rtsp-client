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

	counter := 0
	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error packet")
			continue
		}

		cpy := make([]byte, n)
		copy(cpy, buf)
		go handleBuf(cpy)

		counter++
		if counter > 1 {
			break
		}
	}
}

func handleBuf(buf []byte) {
	/*
		Version := buf[0] & 0x03
		Padding := buf[0]&1<<2 != 0
		Ext := buf[0]&1<<3 != 0
		CSRC := make([]uint, buf[0]>>4)
		Marker := buf[1]&1 != 0
		PayloadType := buf[1] >> 1
		SequenceNumber := toUint(buf[2:4])
		Timestamp := toUint(buf[4:8])
		SyncSource := toUint(buf[8:12])
		fmt.Println(Version, Padding, Ext, CSRC, Marker, PayloadType, SequenceNumber, Timestamp, SyncSource)
	*/
	V := buf[0] >> 4
	P := buf[2]
	X := buf[3]
	CSRS := buf[4:8]

	fmt.Println(V, P, X, CSRS)

	fmt.Println(buf)
}
