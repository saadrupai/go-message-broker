package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saadrupai/go-message-broker/app/config"
	"github.com/saadrupai/go-message-broker/app/container"
	"github.com/saadrupai/go-message-broker/app/server"
)

func main() {
	g := gin.Default()

	config.SetConfig()

	server.StartSocketServer()

	container.Serve(g)

	fmt.Println("Server starting..., pid: ", strconv.Itoa(os.Getpid()))

}
