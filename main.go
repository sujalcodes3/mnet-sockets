package main

import (
	"flag"
	"fmt"
	"github.com/sujalcodes3/media_net_sre_machine_coding/client"
	"github.com/sujalcodes3/media_net_sre_machine_coding/server"
)

func main() {
	machineType := flag.String("type", "server", "specify the type of machine that you want to start up [server, client, heartbeatserver, heartbeatclient]")
	flag.Parse()

	fmt.Printf("[type:%s]\n", *machineType)
	switch *machineType {
	case "server":
		server.StartServer(5000)
	case "heartbeatserver":
		server.StartHeartBeatServer(5000)
	case "commanddispatcherserver":
		server.CommandDispatcherServer(5000)
	case "client":
		client.StartClient(5000)
	case "heartbeatclient":
		client.StartHeartBeatClient(5000)
	case "commanddispatcherclient":
		client.CommandExecutorClient(5000)
	}
}
