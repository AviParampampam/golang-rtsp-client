package main

import (
	"fmt"

	"./rtsp"
)

func main() {
	client := rtsp.NewClient()

	session, err := client.NewSession("192.168.1.15:554", "12345678")
	if err != nil {
		panic(err)
	}
	defer session.Disconnect()

	res, err := session.Setup("rtsp://admin:12345678@192.168.1.15:554/ch01.264?dev=1")
	if err != nil {
		panic(err)
	}

	fmt.Println(res.String())
}
