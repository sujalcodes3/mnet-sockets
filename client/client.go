package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func StartClient(serverPort uint16) {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		panic(fmt.Sprintf("client:while resolving ip addr:%s", err.Error()))
	}

	conn, err := net.DialUDP("udp4", nil, s)
	defer conn.Close()
	if err != nil {
		panic(fmt.Sprintf("client:while dialing the server:%s", err.Error()))
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("client:toserver:")
	str, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("client:while reading from stdin:%s", err.Error()))
	}

	_, er := conn.Write([]byte(str))
	if er != nil {
		panic(fmt.Sprintf("client:while writing to the server:%s", er.Error()))
	}

	buffer := make([]byte, 1024)
	_, _, e := conn.ReadFromUDP(buffer)

	if e != nil {
		panic(fmt.Sprintf("client:while reading from the server"))
	}
	fmt.Printf("client:fromserver:%s", string(buffer))
}

func StartHeartBeatClient(serverPort uint16) {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		panic(fmt.Sprintf("client:while resolving ip addr:%s", err.Error()))
	}

	conn, err := net.DialUDP("udp4", nil, s)
	defer conn.Close()
	if err != nil {
		panic(fmt.Sprintf("client:while dialing the server:%s", err.Error()))
	}

	fmt.Printf("client:specify heartbeat tag:")

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("client:toserver:")
	tag, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("client:while reading from stdin:%s", err.Error()))
	}

	for {
		fmt.Printf("client:heartbeating(%s)...\n", strings.Trim(tag, "\n"))
		_, er := conn.Write([]byte(strings.Trim(tag, "\n")))
		if er != nil {
			panic(fmt.Sprintf("client:while writing to the server:%s", er.Error()))
		}

		buffer := make([]byte, 1024)
		_, _, e := conn.ReadFromUDP(buffer)

		if e != nil {
			panic(fmt.Sprintf("client:while reading from the server"))
		}
		fmt.Printf("client:fromserver:%s\n", string(buffer[:len(buffer)-1]))

		fmt.Printf("client:sleep 1 sec...\n")
		time.Sleep(1 * time.Second)
	}
}

func HeartBeat(conn *net.UDPConn, tag *string) { 
	for {
		fmt.Printf("client:heartbeating(%s)...\n", strings.Trim(*tag, "\n"))
		_, er := conn.Write([]byte(strings.Trim(*tag, "\n")))
		if er != nil {
			panic(fmt.Sprintf("client:while writing to the server:%s", er.Error()))
		}

		buffer := make([]byte, 1024)
		_, _, e := conn.ReadFromUDP(buffer)

		if e != nil {
			panic(fmt.Sprintf("client:while reading from the server"))
		}
		cmd := strings.Trim(string(buffer[:len(buffer)-1]), "{")
		cmd = strings.Trim(string(buffer[:len(buffer)-1]), "}")
		cmd = strings.TrimLeft(string(buffer[:len(buffer)-1]), " ")
		cmd = strings.TrimRight(string(buffer[:len(buffer)-1]), " ")
		fmt.Printf("client:fromserver:%s\n", cmd)

		fmt.Printf("client:sleep 1 sec...\n")
		time.Sleep(1 * time.Second)
	}
}

func CommandExecutorClient(port uint16) {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("client:while resolving ip addr:%s", err.Error()))
	}

	conn, err := net.DialUDP("udp4", nil, s)
	defer conn.Close()
	if err != nil {
		panic(fmt.Sprintf("client:while dialing the server:%s", err.Error()))
	}

	fmt.Printf("client:specify heartbeat tag:")

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("client:toserver:")
	tag, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("client:while reading from stdin:%s", err.Error()))
	}

	go HeartBeat(conn, &tag)

	for { 
		fmt.Printf("client:change tag:\n")
		tg, er := reader.ReadString('\n')
		tag = tg
		if er != nil {
			panic(fmt.Sprintf("client:while reading from stdin"))
		}
	}
}
