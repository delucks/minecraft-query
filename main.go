/*
	How do we want to use this library?

	server = MinecraftServer{IP,
*/
package main

import (
	"fmt"
	"net"
	"os"
)

type MinecraftServer struct {
	connection net.UDPConn
}

func (s MinecraftServer) Connect() {
}

func WriteCommand(sock *net.UDPConn, command byte, body []byte) ([]byte, error) {
	fullPayload := []byte{0xfe, 0xfd, command, 0x01, 0x02, 0x03, 0x04}
	if body != nil {
		fullPayload = append(fullPayload, body...)
	}
	_, err := sock.Write(fullPayload)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 4096)
	nBytes, _, err := sock.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}
	response := buffer[0:nBytes]
	fmt.Printf("Reply: %v\n", response)
	return response, nil
}

func main() {
	host := "localhost"
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, 25565))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	errorPipe := make(chan error, 1)
	go func() {
		// Handshake
		_, err := WriteCommand(conn, 0x09, nil)
		if err != nil {
			errorPipe <- err
			return
		}
	}()

	select {
	case err = <-errorPipe:
		fmt.Println(err)
		os.Exit(1)
	}
}
