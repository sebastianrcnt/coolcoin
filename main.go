package main

import "coolcoin/server"

func main() {
	server := server.CreateServer()
	server.Run("localhost:3000")
}
