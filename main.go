package main

import (
	"fmt"
	"time"

	"./rtp"
	"./rtsp"
)

func main() {
	var res rtsp.Response
	var err error

	listenUDPServers()
	time.Sleep(time.Second * 2)
	client := rtsp.NewClient()

	session, err := client.NewSession(
		"192.168.1.87:554",
		"rtsp://admin:12345678@192.168.1.87:554//h264Preview_01_main",
		"1352; timeout=60",
		"admin",
		"12345678",
	)
	if err != nil {
		panic(err)
	}
	defer session.Disconnect()
	fmt.Println("Session created")

	// OPTIONS
	// res, err = session.Options()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.String())

	// DESCRIBE
	res, err = session.Describe()
	if err != nil {
		panic(err)
	}
	fmt.Println(res.String())

	// SETUP
	// res, err = session.Setup()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.String())

	// PLAY
	// res, err = session.Play()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.String())

	// SETUP AND PLAY
	res, err = session.SetupPlay()
	if err != nil {
		panic(err)
	}
	fmt.Println(res.String())
	time.Sleep(time.Second * 10)
}

func listenUDPServers() {
	go rtp.Listen(41770)
	go rtp.Listen(41771)
}
