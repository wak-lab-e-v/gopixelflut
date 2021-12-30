// gopixelflut
// Package main is the entry-point for the go-sockets server sub-project.
// The go-sockets project is available under the GPL-3.0 License in LICENSE.
package main

import (
	//	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// Application constants, defining host, port, and protocol.
const (
	connHost  = "0.0.0.0"
	connPort  = "1337"
	connType  = "tcp"
	display_x = 60
	display_y = 33
	valueMax  = 7
	debug     = false
)

/* global variable declaration */
var matrix [display_x][display_y]string
var matrix_rgb [display_x][display_y][3]byte

func main() {

	// init display
	fmt.Printf("Init display with %v x %v\n", display_x, display_y)

	for i := 0; i < display_x; i++ {
		for j := 0; j < display_y; j++ {

			matrix[i][j] = "000000"

			matrix_rgb[i][j][0] = 0
			matrix_rgb[i][j][1] = 0
			matrix_rgb[i][j][2] = 0
		}
	}

	// Start the server and listen for incoming connections.
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	// run loop forever, until exit.
	for {
		// Listen for an incoming connection.
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}

		fmt.Println("Client connected.")

		// Print client connection address.
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		// Handle connections concurrently in a new goroutine.
		go handleConnection(c)
	}
}

// handleConnection handles logic for a single connection request.
func handleConnection(conn net.Conn) {

	defer conn.Close()
	bufOne := make([]byte, 1)
	command := ""

	// conn.SetReadDeadline(5) // read timeout

	for {

		n, err := conn.Read(bufOne)

		if err != nil {
			// errors.Is()
			// possible reason: read timeout
			log.Println("error in connection - closed")
			// conn.Close()

			return
		}

		if n > 0 {

			ch := int(bufOne[0])

			if ch == 10 { // ende vom eingehenden command
				resultString := handleCommand(command, conn)
				conn.Write([]byte(resultString + "\r\n"))
				command = ""
			} else {
				command = command + string(bufOne)
			}
		} else {
			log.Println("conn read (n) has 0 bytes")
		}

	}
}

func handleCommand(Command string, conn net.Conn) string {

	if len(Command) < 2 {
		return "cmd len to small"
	}

	CMD := Command[0:2]

	switch CMD {
	case "PX":
		if len(Command) < 5 {
			return "wrong cmd len"
		}
		ARG := Command[3:]
		xyc := strings.Split(ARG, " ")
		if len(xyc) < 3 {
			return "Too few arguments."
		}
		if debug == true {
			log.Println("DEBUG: Full IN: ", xyc)
			log.Println("DEBUG: SP Data X: ", xyc[0])
			log.Println("DEBUG: SP Data Y: ", xyc[1])
			log.Println("DEBUG: SP Data C: ", xyc[2])
		}
		// convert x to int
		xInt, err := strconv.Atoi(xyc[0])

		if err != nil {
			return "Error in X."
		}

		if xInt > display_x {
			return "X to big."
		}

		if xInt == 0 {
			return "X to small."
		}

		// convert y to int
		yInt, err := strconv.Atoi(xyc[1])

		if err != nil {
			return "Error in Y."
		}

		if yInt > display_y {
			return "Y to big."
		}

		if yInt == 0 {
			return "Y to small."
		}

		xyc[2] = strings.TrimRight(xyc[2], "\r\n")

		if len(xyc[2]) != 7 {
			return "Value size missmatch."
		}

		matrix[xInt-1][yInt-1] = xyc[2][1:7]
		//              log.Println("SP from " + xyc[0] + "x" + xyc[1] + " to " + xyc[2] + " from " + conn.RemoteAddr().String())

		return "PX" + xyc[0] + " " + xyc[1] + " " + xyc[2]

	// case "GP":
	// 	// Get Pixel
	// 	if len(Command) < 5 {
	// 		return
	// 	}
	// 	ARG := Command[3:]
	// 	xy := strings.Split(ARG, " ")

	// 	if len(xy) < 2 {
	// 		return "Too few arguments."
	// 	}

	// 	// convert x to int
	// 	xInt, err := strconv.Atoi(xy[0])

	// 	if err != nil {

	// 		return "Error in X."
	// 	}

	// 	if xInt > display_x {
	// 		conn.Write([]byte("X to big."))
	// 		return
	// 	}

	// 	if xInt == 0 {
	// 		conn.Write([]byte("X to small."))
	// 		return
	// 	}

	// 	xy[1] = strings.TrimRight(xy[1], "\r\n")
	// 	// convert y to int
	// 	yInt, err := strconv.Atoi(xy[1])

	// 	if err != nil {
	// 		e := fmt.Errorf("%v", err)
	// 		conn.Write([]byte("Error in Y." + string(e.Error())))
	// 		return
	// 	}

	// 	if yInt > display_y {
	// 		conn.Write([]byte("Y to big."))
	// 		return
	// 	}

	// 	if yInt == 0 {
	// 		conn.Write([]byte("Y to small."))
	// 		return
	// 	}

	// 	conn.Write([]byte("#" + matrix[xInt-1][yInt-1] + "\r\n"))

	// 	return
	// 	//log.Print("GP from " + conn.RemoteAddr().String())

	case "GM":
		// Get Matrix

		m := ""
		for j := 0; j < display_y; j++ {
			for i := 0; i < display_x; i++ {
				m = m + matrix[i][j]
			}
		}
		log.Print("GM from " + conn.RemoteAddr().String() + " Len " + string(len(m)))
		return m

	default:
		if []byte(CMD)[0] > 0 {
			return "unkown command"
		}

	}

	return ""

}
