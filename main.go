package main

import (
	"fmt"

	"./rtsp"
)

const rtspURL = "rtsp://192.168.1.87:554/h264Preview_01_main"

func main() {
	var res rtsp.Response
	var err error

	client := rtsp.NewClient()

	session, err := client.NewSession(
		"192.168.1.87:554",
		"rtsp://192.168.1.87:554/h264Preview_01_main",
		"0056",
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
	// res, err = session.Describe()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.String())

	// SETUP
	res, err = session.Setup()
	if err != nil {
		panic(err)
	}
	fmt.Println(res.String())

	// PLAY
	// res, err = session.Play()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.String())
}
