package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. Create a socket (AF_INET = IPv4, SOCK_STREAM = TCP)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalln("Error in syscall.Socket:", err)
	}
	defer func() {
		err = syscall.Close(fd)
		if err != nil {
			log.Fatalln("Error in syscall.Close:", err)
		}
		fmt.Println("socket", fd, "is closed.")
	}()

	// 1-1. Create a channel to receive OS Signals
	ch := make(chan os.Signal, 1)
	go handleSignal(fd, ch)

	// 2. Create a SockaddrInet4 structure to specify server address and port
	sockAddr := &syscall.SockaddrInet4{Port: 8080}
	copy(sockAddr.Addr[:], []byte{0, 0, 0, 0}) // Listen on all interfaces

	// 3. Bind the socket to the specified address and port
	err = syscall.Bind(fd, sockAddr)
	if err != nil {
		log.Fatalln("Error in syscall.Bind:", err)
	}

	// 4. Set the socket to listen for incoming connections
	err = syscall.Listen(fd, 10)
	if err != nil {
		log.Fatalln("Error in syscall.Listen:", err)
	}

	fmt.Println("Listening on ", "http://localhost:8080")

	for {
		// 5. Accept an incoming connection
		clientFd, sockAddr, err := syscall.Accept(fd)
		if err != nil {
			log.Fatalln("Error in syscall.Accept:", err)
		}

		handleConnection(clientFd, sockAddr.(*syscall.SockaddrInet4))
	}
}

func handleConnection(fd int, sockAddr *syscall.SockaddrInet4) {
	defer func() {
		err := syscall.Close(fd)
		if err != nil {
			log.Fatalln("Error in syscall.Close:", err)
		}
		fmt.Println("socket", fd, "is closed.")
	}()

	// 1. Allocate a buffer to hold incoming data
	buf := make([]byte, 1024)

	for {
		// 2. Read data from the client
		n, err := syscall.Read(fd, buf)
		if err != nil {
			log.Fatalln("Error in syscall.Read:", err)
		}

		// 3. Check for EOF (client closed connection)
		if n == 0 {
			return
		}

		// 4. Determine client IP address
		clientIP := fmt.Sprintf("%d.%d.%d.%d", sockAddr.Addr[0], sockAddr.Addr[1], sockAddr.Addr[2], sockAddr.Addr[3])

		// 5. Print received message along with client IP
		fmt.Printf("Received: %s\nFrom %s\n\n---\n", string(buf[:n]), clientIP)

		// 6. Echo the received message back to the client
		_, err = syscall.Write(fd, buf[:n])
		if err != nil {
			log.Fatalln("Error in syscall.Write:", err)
		}
	}
}

func handleSignal(fd int, ch chan os.Signal) {
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM) // Notify SIGINT and SIGTERM signals

	sig := <-ch
	fmt.Println("\nReceived signal:", sig)

	err := syscall.Close(fd)
	if err != nil {
		log.Fatalln("Error in syscall.Close:", err)
	}

	fmt.Println("socket", fd, "is closed.")
	os.Exit(0)
}
