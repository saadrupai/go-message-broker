package server

import (
	"log"
	"net"

	"github.com/saadrupai/go-message-broker/app/config"
)

func StartSocketServer() {
	port := ":" + config.LocalConfig.SockerServerPort

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to start socket server")
		return
	}

	config.LocalConfig.Listener = listener

}
