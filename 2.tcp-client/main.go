package main

import (
	"fmt"
	"log"
	"syscall"
)

func main() {
	// 1. Create a socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalln("Error in syscall.Socket:", err)
	}
	defer func(fd int) {
		err = syscall.Close(fd)
		if err != nil {
			log.Fatalln("Error in syscall.Close:", err)
		}
	}(fd)

	// 2. Create a SockaddrInet4 structure to specify server address and port
	serverAddr := &syscall.SockaddrInet4{Port: 8080}
	copy(serverAddr.Addr[:], []byte{127, 0, 0, 1}) // 127.0.0.1 (localhost)

	// 3. Connect to the server
	err = syscall.Connect(fd, serverAddr)
	if err != nil {
		log.Fatalln("Error in syscall.Connect:", err)
	}

	// 4. Prompt user for input
	fmt.Print("Enter message: ")
	var message string
	_, err = fmt.Scanln(&message)
	if err != nil {
		log.Fatalln("Error in fmt.Scanln", err)
	}

	// 5. Send message to server
	_, err = syscall.Write(fd, []byte(message))
	if err != nil {
		log.Fatalln("Error in syscall.Write:", err)
	}

	// 6. Read response from server
	buf := make([]byte, 1024)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		log.Fatalln("Error in syscall.Read:", err)
	}

	// 7. Print server response
	fmt.Printf("Server: %s\n", string(buf[:n]))
}
