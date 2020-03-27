/*
	How do we want to use this library?

	server = MinecraftServer{IP,
*/
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
)

type MinecraftServer struct {
	connection net.UDPConn
}

func (s MinecraftServer) Connect() {
}

func WriteCommand(sock *net.UDPConn, command byte, body []byte) ([]byte, error) {
	fullPayload := []byte{0xfe, 0xfd, command, 0x00, 0x00, 0x00, 0x01}
	if body != nil {
		fullPayload = append(fullPayload, body...)
	}
	fmt.Printf("Request: %v\n", fullPayload)
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
	if len(response) < 5 || response[0] != command {
		return nil, fmt.Errorf("Response %v was too short or didn't match with command %v\n", response, command)
	}
	fmt.Printf("Reply: %v\n", response)
	return response[5:], nil
}

func swapByteOrder(in []byte) []byte {
	for i := len(in)/2 - 1; i >= 0; i-- {
		opp := len(in) - 1 - i
		in[i], in[opp] = in[opp], in[i]
	}
	return in
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
		response, err := WriteCommand(conn, 0x09, nil)
		if err != nil {
			errorPipe <- err
			return
		}
		// Parse challenge token from the response
		nonNullTerminated := string(response[:len(response)-1])
		parsedToken, err := strconv.ParseInt(nonNullTerminated, 10, 32)
		if err != nil {
			errorPipe <- err
			return
		}
		buf := new(bytes.Buffer)
		// convert to int32 to correct the binary size
		err = binary.Write(buf, binary.BigEndian, int32(parsedToken))
		if err != nil {
			errorPipe <- err
			return
		}
		packedToken := buf.Bytes()
		// Must be padded with four null bytes
		challengeResponse := append(packedToken, []byte{0x00, 0x00, 0x00, 0x00}...)
		_, err = WriteCommand(conn, 0x00, challengeResponse)
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
