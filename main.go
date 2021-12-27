// based on https://github.com/Alice-Williams-Tech/go-sockets
// This and the go-sockets project is available under the GPL-3.0 License in LICENSE.
//
package main

import (
	"bufio"
	//	"encoding/hex"
	"errors"
	"fmt"
	"image/color"
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
var errInvalidFormat = errors.New("invalid hex color format")

func main() {

	//
	// init display/canvas/matrix
	//
	fmt.Printf("Init display with %v x %v\n", display_x, display_y)
	for i := 0; i < display_x; i++ {
		for j := 0; j < display_y; j++ {

			matrix[i][j] = "000000"

			matrix_rgb[i][j][0] = 0
			matrix_rgb[i][j][1] = 0
			matrix_rgb[i][j][2] = 0
		}
	}

	//
	// Start the server and listen for incoming connections.
	//
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

func ParseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
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

	switch command := string(buffer[0:2]); command {
	case "IH":
	case "CC":
	case "SP":
	case "GP":
	case "GM":
		getMatrix(conn)
	}

	//
	// Send some info on request
	//
	if string(buffer[0:1]) == "I" {

	}

	//
	// Command & Control (CC)
	//
	if string(buffer[0:1]) == "CC" {

	}

	//
	// Set Pixel (SP)
	//
	// Sample: 	SP X1 Y1 #FFFFFF
	//			SP X01 Y33 R0 G255 B17     <-- faster
	//

	if string(buffer[0:2]) == "SP" {

		xyc := strings.Split(string(buffer[3:]), " ")

		if debug == true {
			log.Println("DEBUG: Set Pixel: ", xyc)
		}

		if (len(xyc) < 3) || (len(xyc) > 5) {
			conn.Write([]byte("Too few or many arguments. Use I command."))
			conn.Close()
			return
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

		xyc[2] = strings.TrimRight(xyc[2], "\r\n")

		if len(xyc[2]) != 7 {
			conn.Write([]byte("Value size missmatch."))
			conn.Close()
			return
		}

		// set 3. value to display matrix
		matrix[xInt-1][yInt-1] = xyc[2][1:7]

		matrix_rgb[xInt-1][yInt-1][0] = []byte(xyc[2][1:2])[0]
		// hex.Decode()
		log.Println("SP x:" + xyc[0] + " y:" + xyc[1] + " to " + xyc[2] + " from " + conn.RemoteAddr().String())

		bufferOut = []byte("OK, " + string(buffer[3:]) + "\n")
	}

	//
	// Get Pixel
	// Get color from X/Y
	// Result as #FFFFF
	//
	if string(buffer[0:2]) == "GP" {

		xy := strings.Split(string(buffer[3:]), " ")

		if len(xy) < 2 {
			conn.Write([]byte("Too few arguments."))
			conn.Close()
			return
		}

		// convert x to int
		xInt, err := strconv.Atoi(xy[0])

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

		xy[1] = strings.TrimRight(xy[1], "\r\n")
		// convert y to int
		yInt, err := strconv.Atoi(xy[1])

		if err != nil {
			e := fmt.Errorf("%v", err)
			conn.Write([]byte("Error in Y." + string(e.Error())))
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

		bufferOut = []byte("#" + matrix[xInt-1][yInt-1] + "\r\n")

		log.Print("GP from " + conn.RemoteAddr().String())
	}

	// Print response message, stripping newline character.
	// log.Println("Client message:", string(buffer[:len(buffer)-1]))

	// Send response message to the client.
	conn.Write(bufferOut)

	// Restart the process.
	handleConnection(conn)
}

func infoAndHelp() {

}

//
// Get Matrix
// Send full matrix row 1, col 1 to X, row 2, col 1 to X and so on
// Result as stream of ASCII endoded Hex (RGB)
//
func getMatrix(conn net.Conn) {
	var checksum uint64 = 0

	for j := 0; j < display_y; j++ {
		for i := 0; i < display_x; i++ {
			conn.Write([]byte(matrix[i][j]))

			check, err := strconv.Atoi(matrix[i][j]) // build checksum
			if err == nil {
				checksum = checksum + uint64(check)
			}

		}
	}

	check_string := strconv.FormatUint(checksum, 10)
	conn.Write([]byte("\r\n#" + check_string))
	conn.Write([]byte("\r\n"))

	log.Println("GM from " + conn.RemoteAddr().String())

}
