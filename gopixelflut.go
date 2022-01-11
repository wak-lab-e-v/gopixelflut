// gopixelflut
// Package main is the entry-point for the go-sockets server sub-project.
// The go-sockets project is available under the GPL-3.0 License in LICENSE.
package main

import (
	//	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Application constants, defining host, port, and protocol.
const (
	connHost  = "0.0.0.0"
	connPort  = "1337"
	connType  = "tcp"
	display_y = 33
	display_x = 60
)

/* global variable declaration */

var serverStartTime time.Time
var debug bool = false
var matrix_rgb [display_x][display_y][3]int64
var connectionCount uint64 = 0
var activeConnCount uint64 = 0
var commandCount uint64 = 0
var errorCount uint64 = 0

func main() {

	for _, arg := range os.Args[1:] {
		if arg == "debug" {
			debug = true
			fmt.Printf("DEBUG activated!\n")
		}

	}

	serverStartTime = time.Now()

	// init display
	fmt.Println("PixlServer - init display with %v x %v\n", display_x, display_y)

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
			errorCount++
			return
		}

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		// Handle connections concurrently in a new goroutine.
		go handleConnection(c)

	}
}

// handleConnection handles logic for a single connection request.
func handleConnection(conn net.Conn) {

	connectionCount++
	activeConnCount++

	bufOne := make([]byte, 1)
	command := ""
	connTimeOut := 0

	// conn.SetReadDeadline(5) // read timeout

	for {
		// defer conn.Close()

		n, err := conn.Read(bufOne)

		if err != nil {
			if errors.Is(err, io.EOF) {
				time.Sleep(50 * time.Microsecond)
				connTimeOut++
				if connTimeOut > 1000 {
					fmt.Errorf("Timeout - Conn closed")
					conn.Write([]byte("Timeout - Conn closed"))
					conn.Close()
					decrementActiveConnectionCounter()
					return
				}
				continue

			} else {
				fmt.Errorf("Error: %v", err)
				conn.Write([]byte(err.Error() + "\r\n"))
				conn.Close()
				decrementActiveConnectionCounter()
				return
			}
		}

		if n > 0 {

			if bufOne[0] == 10 { // end incomming command
				resultString := handleCommand(command, conn)
				conn.Write([]byte(resultString + "\r\n"))
				command = ""

			} else {
				command = command + string(bufOne)
			}
		}
	}
}

func handleCommand(Command string, conn net.Conn) string {

	commandCount++

	if debug {
		log.Println("DEBUG: "+conn.RemoteAddr().String()+" command len: ", len(Command))
		log.Println("DEBUG: "+conn.RemoteAddr().String()+" Full command: ", Command)
	}

	if len(Command) < 2 {
		return "Unkown command. Use 'HELP' and 'INFO'"
	}

	if Command[0:2] == "PX" { // SET or GET Pixel

		if debug {
			log.Println("PX")
		}

		ARG := strings.Split(Command[3:], " ")

		// PX X Y #000000 (ASCI-hex encoded RGB value), ARG len = 3
		// PX X Y 0 0 0  (ASCI-int based RGB value) , ARG len = 5
		// ARG[0] = X
		// ARG[1] = Y
		// ARG[2] = #00000 or 0
		// ARG[3] = 0
		// ARG[4] = 0

		//
		// Start with check of incomming X/Y
		//

		// convert x to int
		xInt, err := strconv.Atoi(ARG[0])
		if err != nil {
			errorCount++
			return "Error in X. Use 'HELP'"
		}
		if xInt > display_x {
			errorCount++
			return "X to big. Use 'INFO'."
		}
		if xInt == 0 {
			errorCount++
			return "X to small. Use 'INFO'."
		}

		// convert y to int
		yInt, err := strconv.Atoi(ARG[1])
		if err != nil {
			errorCount++
			return "Error in Y. Use 'HELP'" + err.Error()
		}
		if yInt > display_y {
			errorCount++
			return "Y to big. Use 'INFO'"
		}
		if yInt == 0 {
			errorCount++
			return "Y to small. Use 'INFO'"
		}

		//
		// Start selecting SET / GET Modes
		//
		if debug {
			log.Println("Start selecting SET / GET Modes")
		}

		if len(ARG) == 5 { // MODE: SET PIXEL COLOR by 0 0 0 - 255 255 255
			if debug {
				log.Println("Detect Mode: PX X Y R G B")
				log.Println("ARG 0 = X = " + ARG[0])
				log.Println("ARG 1 = Y = " + ARG[1])
				log.Println("ARG 2 = R = " + ARG[2])
				log.Println("ARG 3 = G = " + ARG[3])
				log.Println("ARG 4 = B = " + ARG[4])
			}

			r, _ := strconv.Atoi(ARG[2])
			g, _ := strconv.Atoi(ARG[3])
			b, _ := strconv.Atoi(strings.TrimRight(ARG[4], "\r\n"))

			matrix_rgb[xInt-1][yInt-1][0] = int64(r)
			matrix_rgb[xInt-1][yInt-1][1] = int64(g)
			matrix_rgb[xInt-1][yInt-1][2] = int64(b)
			return Command
		}

		if len(ARG) == 3 { // MODE: SET PIXEL COLOR by #000000 - #FFFFFF

			if debug {
				log.Println("Detect Mode: PX X Y #RRGGBB")
				log.Println("ARG 0 = X = " + ARG[0])
				log.Println("ARG 1 = Y = " + ARG[1])
				log.Println("ARG 2 = C = " + ARG[2])
			}

			if len(ARG[2]) < 7 {
				errorCount++
				return "Error in color syntax. Use 'HELP'"
			}

			r, _ := strconv.ParseInt(ARG[2][1:3], 16, 64) // hex to int64
			g, _ := strconv.ParseInt(ARG[2][3:5], 16, 64)
			b, _ := strconv.ParseInt(ARG[2][5:7], 16, 64)

			matrix_rgb[xInt-1][yInt-1][0] = r
			matrix_rgb[xInt-1][yInt-1][1] = g
			matrix_rgb[xInt-1][yInt-1][2] = b

			return Command
		}

		errorCount++
		return "Count of arguments not valid. Use 'HELP'"

	}

	if Command[0:2] == "GM" { // Get full pixelmatrix

		var m []byte

		for j := 0; j < display_y; j++ {
			for i := 0; i < display_x; i++ {
				m = append(m, byte(matrix_rgb[i][j][0]))
				m = append(m, byte(matrix_rgb[i][j][1]))
				m = append(m, byte(matrix_rgb[i][j][2]))
			}
		}
		conn.Write(m)
		if debug == true {
			log.Print("GM from " + conn.RemoteAddr().String())
		}
		return ""
	}

	if Command[0:2] == "GP" {

		ARG := strings.Split(Command[3:], " ")
		// convert x to int
		xInt, err := strconv.Atoi(ARG[0])
		if err != nil {
			errorCount++
			return "Error in X. Use 'HELP'" + err.Error()
		}
		if xInt > display_x {
			errorCount++
			return "X to big. Use 'INFO'."
		}
		if xInt == 0 {
			errorCount++
			return "X to small. Use 'INFO'."
		}

		ARG[1] = strings.TrimRight(ARG[1], "\r\n")

		// convert y to int
		yInt, err := strconv.Atoi(ARG[1])
		if err != nil {
			errorCount++
			return "Error in Y. Use 'HELP'" + err.Error()
		}
		if yInt > display_y {
			errorCount++
			return "Y to big. Use 'INFO'"
		}
		if yInt == 0 {
			errorCount++
			return "Y to small. Use 'INFO'"
		}

		r := strconv.Itoa(int(matrix_rgb[xInt-1][yInt-1][0]))
		g := strconv.Itoa(int(matrix_rgb[xInt-1][yInt-1][1]))
		b := strconv.Itoa(int(matrix_rgb[xInt-1][yInt-1][2]))

		return "PX " + ARG[0] + " " + ARG[1] + " " + r + " " + g + " " + b
	}

	if len(Command) < 4 {
		errorCount++
		return "Unkown command. Use 'HELP' and 'INFO'"
	}

	if Command[0:4] == "HELP" { // Small HELP
		return "PixelServer in GO by der-den (WAK-Lab). See https://github.com/wak-lab-e-v/gopixelflut"
	}

	if Command[0:4] == "INFO" { // Some Infos
		infoX := strconv.Itoa(display_x)
		infoY := strconv.Itoa(display_y)
		info := "Matrix Size: " + infoX + "x" + infoY + "\r\n"
		info = info + "Server start: " + serverStartTime.Format(time.RFC3339) + "\r\n"
		diff := time.Now().Sub(serverStartTime)
		info = info + "Runtime: " + diff.String() + "\r\n"
		info = info + "Active connections: " + strconv.FormatUint(activeConnCount, 10) + "\r\n"
		info = info + "Connection counter: " + strconv.FormatUint(connectionCount, 10) + "\r\n"
		info = info + "Command counter: " + strconv.FormatUint(commandCount, 10) + "\r\n"
		info = info + "Error counter: " + strconv.FormatUint(errorCount, 10) + "\r\n"
		return info
	}

	if Command[0:4] == "SIZE" {
		infoX := strconv.Itoa(display_x)
		infoY := strconv.Itoa(display_y)
		return "SIZE: " + infoX + " " + infoY
	}

	if Command[0:4] == "EXIT" { // Exit connection in
		conn.Close()
		// decrementActiveConnectionCounter()
		return ""
	}

	return "Unkown command. Use 'HELP' and 'INFO'"

}

func decrementActiveConnectionCounter() {
	if activeConnCount > 0 {
		activeConnCount--
	}
}
