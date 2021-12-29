// Package main is the entry-point for the go-sockets server sub-project.
// The go-sockets project is available under the GPL-3.0 License in LICENSE.
package main

import (
	"bufio"
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

func main() {

	// init display
	fmt.Printf("Init display with %v x %v\n", display_x, display_y)

	for i := 0; i < display_x; i++ {
		for j := 0; j < display_y; j++ {

			matrix[i][j] = "000000"
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
	//buffer := make([]byte, 1024)

	for {
		// Buffer client input until a newline.
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')

		//bufLen , err := conn.Read(buffer)
		//bufferstring := string(buffer)
		//log.Println(bufLen)
		//log.Println(bufferstring)
		//buffersplit := strings.Split(string(buffer), "\n")
		//bufferstring = buffersplit[0]
		//log.Println(bufferstring)
		//bufferstring := string(buffer)
		// log.Println(bufferstring)
		// Close left clients.
		if err != nil {
			conn.Close()
			return
		}
		go handleCommand(buffer, conn)
	}
}

func handleCommand(Command []byte, conn net.Conn) {
	log.Println(string(Command))

	if len(Command) < 3 {
		return
	}

	CmdString := string(Command)
	CMD := CmdString[0:2]
	ARG := CmdString[3:]

	switch CMD {
	case "SP":
		xyc := strings.Split(ARG, " ")
		if len(xyc) < 3 {
			conn.Write([]byte("Too few arguments."))
			return
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
			conn.Write([]byte("Error in X."))
			return
		}

		if xInt > display_x {
			conn.Write([]byte("X to big."))
			return
		}

		if xInt == 0 {
			conn.Write([]byte("X to small."))
			return
		}

		// convert y to int
		yInt, err := strconv.Atoi(xyc[1])

		if err != nil {
			conn.Write([]byte("Error in Y."))
			return
		}

		if yInt > display_y {
			conn.Write([]byte("Y to big."))
			return
		}

		if yInt == 0 {
			conn.Write([]byte("Y to small."))
			return
		}

		xyc[2] = strings.TrimRight(xyc[2], "\r\n")

		if len(xyc[2]) != 7 {
			conn.Write([]byte("Value size missmatch."))
			return
		}

		// set 3. value to display matrix
		matrix[xInt-1][yInt-1] = xyc[2][1:7]
		//              log.Println("SP from " + xyc[0] + "x" + xyc[1] + " to " + xyc[2] + " from " + conn.RemoteAddr().String())

		conn.Write([]byte("OK\n"))
		return

	case "GP":
		// Get Pixel
		xy := strings.Split(ARG, " ")

		if len(xy) < 2 {
			conn.Write([]byte("Too few arguments."))
			return
		}

		// convert x to int
		xInt, err := strconv.Atoi(xy[0])

		if err != nil {
			conn.Write([]byte("Error in X."))
			return
		}

		if xInt > display_x {
			conn.Write([]byte("X to big."))
			return
		}

		if xInt == 0 {
			conn.Write([]byte("X to small."))
			return
		}

		xy[1] = strings.TrimRight(xy[1], "\r\n")
		// convert y to int
		yInt, err := strconv.Atoi(xy[1])

		if err != nil {
			e := fmt.Errorf("%v", err)
			conn.Write([]byte("Error in Y." + string(e.Error())))
			return
		}

		if yInt > display_y {
			conn.Write([]byte("Y to big."))
			return
		}

		if yInt == 0 {
			conn.Write([]byte("Y to small."))
			return
		}

		conn.Write([]byte("#" + matrix[xInt-1][yInt-1] + "\r\n"))

		log.Print("GP from " + conn.RemoteAddr().String())
	case "GM":
		// Get Matrix
		for j := 0; j < display_y; j++ {
			for i := 0; i < display_x; i++ {
				conn.Write([]byte(matrix[i][j]))
			}
		}

		conn.Write([]byte("\r\n"))
		log.Println("GM from " + conn.RemoteAddr().String())
		return

	default:
		conn.Write([]byte("unkown command \n"))
		return
	}

	// Print response message, stripping newline character.
	// log.Println("Client message:", string(buffer[:len(buffer)-1]))

	// Send response message to the client.

}
