package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	//"strconv"
	"sync"
)

func StartServer(port uint16) {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("server:while resolving UDP addr: %s", err.Error()))
	}
	conn, err := net.ListenUDP("udp4", s)
	if err != nil {
		panic(fmt.Sprintf("server:while listening to UDP addr:%s", err.Error()))
	}

	fmt.Printf("server:listening on port %d\n", port)

	defer conn.Close()

	buffer := make([]byte, 1024)
	n, clientAddr, err := conn.ReadFromUDP(buffer)

	fmt.Printf("server:message recieved from (%s): { message: %s , len: %d } echoing...\n", clientAddr.String(), string(buffer), n)
	_, er := conn.WriteToUDP(buffer, clientAddr)

	if er != nil {
		panic(fmt.Sprintf("server:while echoing:%s", err.Error()))
	}
}

func StartHeartBeatServer(port uint16) {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("server:while resolving UDP addr: %s", err.Error()))
	}
	conn, err := net.ListenUDP("udp4", s)
	if err != nil {
		panic(fmt.Sprintf("server:while listening to UDP addr:%s", err.Error()))
	}

	fmt.Printf("server:listening on port %d\n", port)

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buffer)

		fmt.Printf("server:heartbeat request recieved from (%s): { message: %s , len: %d } responding...\n", clientAddr.String(), string(buffer[:len(buffer)-1]), n)
		toBeSent := fmt.Sprintf("{ alive: tag(%s) }\n", string(buffer[len(buffer)-1]))
		_, er := conn.WriteToUDP([]byte(toBeSent), clientAddr)

		if er != nil {
			panic(fmt.Sprintf("server:while echoing:%s", err.Error()))
		}
	}
}

func ParseInt(str string) uint8 {
	str = strings.Trim(str, "\x00")
	if len(str) == 1 {
		key, _ := strconv.Atoi(str)
		return uint8(key)
	}
	var res uint8 = 0
	var mult uint8 = uint8(len(str)-1) * 10
	i := 0
	for i < len(str) {
		key, err := strconv.Atoi((string(str[i])))
		if str[i] == '\x00' { 
			break
		}
		if err != nil {
			fmt.Printf("problem here:%s", err.Error())
		}
		res += mult * uint8(key)
		mult /= 10
		i++
	}

	return res
}

func ListenHeartBeat(conn *net.UDPConn, commandMap *CommandMap) {
	fmt.Printf("server:listening for heartbeat concurrently...\n")
	for {
		buffer := make([]byte, 4096)
		n, clientAddr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Printf("tag format not correct\n")
		}

		tag := ParseInt(string(buffer[:len(buffer) - 1]))
		fmt.Printf("server:heartbeat request recieved from (%s): { message: %v %T , len: %d } responding...\n", clientAddr.String(), tag, tag, n)

		toBeSent := fmt.Sprintf("{ alive: tag(%s) }\n", string(buffer[len(buffer)-1]))

		commandMap.mut.Lock()
		if len(commandMap.commandMap) == 0 {
			toBeSent = fmt.Sprintf("{ <empty> 0 }")
		}
		key, found := commandMap.commandMap[tag]
		if found == false || key == nil { 
			toBeSent = fmt.Sprintf("{ <no tag> }")
		} else { 
			toBeSent = fmt.Sprintf("{ %s }", key.cmd)
		}
		
		commandMap.mut.Unlock()
	
		_, er := conn.WriteToUDP([]byte(toBeSent), clientAddr)
		commandMap.commandMap[tag] = nil

		if er != nil {
			panic(fmt.Sprintf("server:while echoing:%s", err.Error()))
		}
	}
}

type Command struct {
	tag uint8 
	cmd string
}

func (cmd *Command) String() string {
	return fmt.Sprintf("{ %d | %s }", cmd.tag, cmd.cmd)
}

func DeserializeCommand(cmd string) *Command {
	i := 0

	for i < len(cmd) {
		if cmd[i] == ' ' {
			break
		}
		i++
	}
	command := cmd[i:]
	return &Command{
		tag: ParseInt(cmd[:i]),
		cmd: command,
	}
}

type CommandMap struct {
	commandMap map[uint8]*Command
	mut        sync.Mutex
}

func NewCommandMap() *CommandMap {
	return &CommandMap{
		commandMap: make(map[uint8]*Command),
	}
}

func CommandDispatcherServer(port uint16) {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("server:while resolving UDP addr: %s", err.Error()))
	}
	conn, err := net.ListenUDP("udp4", s)
	if err != nil {
		panic(fmt.Sprintf("server:while listening to UDP addr:%s", err.Error()))
	}

	fmt.Printf("server:listening on port %d\n", port)

	defer conn.Close()

	commandMap := NewCommandMap()

	go ListenHeartBeat(conn, commandMap)

	for {
		fmt.Printf("server:cmd:")
		reader := bufio.NewReader(os.Stdin)
		cmdIn, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("while reading command from stdin:%s", err.Error()))
		}
		command := DeserializeCommand(cmdIn)
		fmt.Printf("%s\n", command)

		commandMap.mut.Lock()
		commandMap.commandMap[command.tag] = command
		commandMap.mut.Unlock()
	}
}
