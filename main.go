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
)

/* global variable declaration */
var matrix [display_x][display_y][3]uint8
var matrix_b [display_x][display_y][3]byte
var matrix_s [display_x][display_y]string

func main() {

	// init display
	fmt.Printf("Init display with %v x %v\n", display_x, display_y)

	for i := 0; i < display_x; i++ {
		for j := 0; j < display_y; j++ {
			matrix[i][j][0] = 0
			matrix[i][j][1] = 0
			matrix[i][j][2] = 0
			matrix_b[i][j][0] = 0
			matrix_b[i][j][1] = 0
			matrix_b[i][j][2] = 0
			matrix_s[i][j] = "000000"
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
	// Buffer client input until a newline.
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	// Close left clients.
	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	//var n int8

	// Set Pixel
	if string(buffer[0:2]) == "SP" {
		xyc := strings.Split(string(buffer[3:]), " ")
		x := xyc[0]
		y := xyc[1]
		colorHex := xyc[2]

		log.Println("SP Data X: ", x)
		log.Println("SP Data Y: ", y)
		log.Println("SP Data C: ", colorHex)

		//	if len(colorhex) <> 4 then die

		xInt, err := strconv.Atoi(x)

		if err != nil {
			conn.Write([]byte("Error in X."))
			conn.Close()
			return
		}

		yInt, err := strconv.Atoi(y)

		if err != nil {
			conn.Write([]byte("Error in Y."))
			conn.Close()
			return
		}

		//	matrix_b[xInt][yInt][0] = []byte(r)[0] // use first byte after convert "string" r to byte array
		matrix_s[xInt][yInt] = colorHex[1:]
		//		matrix[xInt][yInt][1] = c[3:2]
		//		matrix[xInt][yInt][2] = c[5:2]

		log.Println("SP Data: ", xyc)
	}

	// Get Pixel
	if string(buffer[0:2]) == "GP" {
		//xyc := strings.Split(string(buffer[2:]), " ")
		log.Println("#", matrix)
	}

	// Get Matrix
	if string(buffer[0:2]) == "GM" {

	}

	// Print response message, stripping newline character.
	log.Println("Client message:", string(buffer[:len(buffer)-1]))

	// Send response message to the client.
	conn.Write(buffer)

	// Restart the process.
	handleConnection(conn)
}
