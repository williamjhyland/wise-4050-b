package tcp_ip

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func RunTCPIP() {
	// Connect to the device
	conn, err := net.Dial("tcp", "10.1.14.38:80") // Replace 'port' with the correct port number
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Set a deadline for reading. Adjust as needed.
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// Send data to the device
	_, err = conn.Write([]byte("Your data here\n")) // Replace with your data
	if err != nil {
		fmt.Println("Error writing to stream:", err)
		return
	}

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from stream:", err)
		return
	}
	fmt.Println("Response from device:", response)
}
