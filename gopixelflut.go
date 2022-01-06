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
var matrix_rgb [display_x][display_y][3]int64

func main() {

	// init display
	fmt.Printf("Init display with %v x %v\n", display_x, display_y)

	for i := 0; i < display_x; i++ {
		for j := 0; j < display_y; j++ {

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

	bufOne := make([]byte, 1)
	command := ""
	var cmd []byte

	// conn.SetReadDeadline(5) // read timeout

	for {
		// defer conn.Close()

		n, err := conn.Read(bufOne)

		if err != nil {
			fmt.Errorf("Error: %v", err)
			conn.Write([]byte(err.Error() + "\r\n"))
			// if errors.Is(err, ...) { }
			// possible reason: read timeout
			log.Println("error in connection - closed")
			conn.Close()

			return
		}

		if n > 0 {

			ch := int(bufOne[0])

			if ch == 10 { // end incomming command
				resultString := handleCommand(command, conn)
				conn.Write([]byte(resultString + "\r\n"))
				command = ""

			} else {
				command = command + string(bufOne)
				cmd = cmd + bufOne
			}
		}
	}
}

func handleCommand(Command string, conn net.Conn) string {

	if Command[0:2] == "PX" { // SET or GET Pixel

		ARG := strings.Split(Command[3:], " ")

		// PX X Y #000000 (ASCI-hex encoded RGB value), ARG len = 3
		// PX X Y 0 0 0  (ASCI-int based RGB value) , ARG len = 5
		// ARG[0] = X
		// ARG[1] = Y
		// ARG[2] = #00000 or 0
		// ARG[3] = 0
		// ARG[4] = 0

		if debug == true {
			log.Println("DEBUG: Full IN: ", ARG)
		}

		//
		// Start with check of incomming X/Y
		//

		// convert x to int
		xInt, err := strconv.Atoi(ARG[0])
		if err != nil {
			return "Error in X. Use 'HELP'"
		}
		if xInt > display_x {
			return "X to big. Use 'INFO'."
		}
		if xInt == 0 {
			return "X to small. Use 'INFO'."
		}

		// convert y to int
		yInt, err := strconv.Atoi(ARG[1])
		if err != nil {
			return "Error in Y. Use 'HELP'"
		}
		if yInt > display_y {
			return "Y to big. Use 'INFO'"
		}
		if yInt == 0 {
			return "Y to small. Use 'INFO'"
		}

		//
		// Start selecting SET / GET Modes
		//
		if len(ARG) == 5 { // MODE: SET PIXEL COLOR by 0 0 0 - 255 255 255
			r, err := strconv.Atoi(ARG[2])
			g, err := strconv.Atoi(ARG[3])
			b, err := strconv.Atoi(ARG[4])
			matrix_rgb[xInt-1][yInt-1][0] = int64(r)
			matrix_rgb[xInt-1][yInt-1][1] = int64(g)
			matrix_rgb[xInt-1][yInt-1][2] = int64(b)
			return Command
		}

		if len(ARG) == 3 { // MODE: SET PIXEL COLOR by #000000 - #FFFFFF
			r, err := strconv.ParseInt(ARG[2][0:2], 16, 64) // hex to int64
			g, err := strconv.ParseInt(ARG[2][3:2], 16, 64)
			b, err := strconv.ParseInt(ARG[2][5:2], 16, 64)

			matrix_rgb[xInt-1][yInt-1][0] = r
			matrix_rgb[xInt-1][yInt-1][1] = g
			matrix_rgb[xInt-1][yInt-1][2] = b
			return Command
		}

		if len(ARG) == 2 { // MODE: GET PIXEL COLOR
			return "PX X Y " + string(matrix_rgb[xInt-1][yInt-1][0])
		}

		return "Count of arguments not valid. Use 'HELP'"

		xyc[2] = strings.TrimRight(xyc[2], "\r\n")

		//              log.Println("SP from " + xyc[0] + "x" + xyc[1] + " to " + xyc[2] + " from " + conn.RemoteAddr().String())

	}

	if Command[0:2] == "GM" { // Get full pixelmatrix
		m := ""
		c := ""
		for j := 0; j < display_y; j++ {
			for i := 0; i < display_x; i++ {
				c = strconv.FormatInt(matrix_rgb[i][j][0], 16)
				//+matrix[i][j][3]+matrix[i][j][3]
				m = m + c
			}
		}
		if debug == true {
			log.Print("GM from " + conn.RemoteAddr().String() + " Len " + string(len(m)))
		}
		return m
	}

	if Command[0:4] == "HELP" { // Small HELP
		return "PixelServer in GO by WAK-Lab."
	}

	if Command[0:4] == "INFO" { // Some Infos
		info := "Matrix Size: \n"
		info = info + "Run since: "
		return info
	}

	return "unkown command"

}
