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
	connHost  = "localhost"
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

	bufferOut := []byte("unkown command \n")

	// Buffer client input until a newline.
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	// Close left clients.
	if err != nil {
		conn.Close()
		return
	}

	//
	// Set Pixel
	//
	if string(buffer[0:2]) == "SP" {

		xyc := strings.Split(string(buffer[3:]), " ")

		if len(xyc) < 3 {
			conn.Write([]byte("Too few arguments."))
			conn.Close()
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
			conn.Close()
			return
		}

		if xInt > display_x {
			conn.Write([]byte("X to big."))
			conn.Close()
			return
		}

		if xInt == 0 {
			conn.Write([]byte("X to small."))
			conn.Close()
			return
		}

		// convert y to int
		yInt, err := strconv.Atoi(xyc[1])

		if err != nil {
			conn.Write([]byte("Error in Y."))
			conn.Close()
			return
		}

		if yInt > display_y {
			conn.Write([]byte("Y to big."))
			conn.Close()
			return
		}

		if yInt == 0 {
			conn.Write([]byte("Y to small."))
			conn.Close()
			return
		}

		// set 3. value to display matrix
		matrix[xInt-1][yInt-1] = xyc[2][1:]
		log.Println("")

		bufferOut = []byte("OK, " + string(buffer[3:]) + "\n")
	}

	// Get Pixel
	if string(buffer[0:2]) == "GP" {
		//xyc := strings.Split(string(buffer[2:]), " ")
		log.Print(".", matrix)
	}

	// Get Matrix
	if string(buffer[0:2]) == "GM" {
		log.Print("#", matrix)
	}

	//

	// Print response message, stripping newline character.
	// log.Println("Client message:", string(buffer[:len(buffer)-1]))

	// Send response message to the client.
	conn.Write(bufferOut)

	// Restart the process.
	handleConnection(conn)
}
