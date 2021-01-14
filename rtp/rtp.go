package rtp

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var file, _ = os.Create("video.mp4")

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

	// counter := 0
	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error packet")
			continue
		}

		cpy := make([]byte, n)
		copy(cpy, buf)
		// go handleBuf(cpy)
		handleBuf(cpy)

		// counter++
		// if counter > 5 {
		// 	break
		// }
	}
}

// Packet is an RTP packet
type Packet struct {
	Version        byte   // Версия протокола
	Padding        byte   // Дополняется ли пустыми байтами на конце
	Extension      byte   // Есть ли расширения протокола
	CSRCCount      byte   // Количество CSRC-идентификаторов
	Marker         byte   // Маркер особого назначения
	PayloadType    byte   // Формат полезной нагрузки
	SequenceNumber byte   // Порядковый номер
	Timestamp      byte   // Метка времени
	SSRC           byte   // SSRC-идентификатор
	CSRC           []byte // CSRC-идентификаторы
	Payload        []byte
}

func handleBuf(buf []byte) {
	var p Packet

	p.Version = buf[0] >> 6                   // 0-1
	p.Padding = (buf[0] >> 5) & 1             // 2
	p.Extension = (buf[0] >> 5) & 1           // 3
	p.CSRCCount = buf[0] & 0x0f               // 4-7
	p.Marker = buf[1] >> 7                    // 8
	p.PayloadType = buf[1] & 0x7f             // 9-15
	p.SequenceNumber = byteAddition(buf[2:4]) // 16-31
	p.Timestamp = byteAddition(buf[4:8])      // 32
	p.SSRC = byteAddition(buf[8:12])          // 64

	// Parsing CSRC
	lastPosition, step := 12, 4
	for i := 0; i < int(p.CSRCCount); i++ {
		p.CSRC[i] = byteAddition(buf[lastPosition : lastPosition+step+1]) // 96
		lastPosition += step
	}

	// Parsing Payload
	p.Payload = buf[lastPosition:]

	writeData(p.Payload)

	// if p.Extension > 0 {
	// 	// TODO: Добавить обработку заголовков расширения
	// 	fmt.Println(p.CSRCCount * 32)
	// }

	fmt.Println(p)
}

func byteAddition(b []byte) (res byte) {
	for _, b := range b {
		res += b
	}
	return res
}

func writeData(data []byte) {
	file.Write(data)
	fmt.Println("Packet writing")
}

// func bitsToBase2(data []byte) string {
// 	var buf bytes.Buffer
// 	for _, b := range data {
// 		fmt.Fprintf(&buf, "%08b ", b)
// 	}
// 	buf.Truncate(buf.Len() - 1) // To remove extra space
// 	return fmt.Sprintf("<%s>\n", buf.Bytes())
// }
