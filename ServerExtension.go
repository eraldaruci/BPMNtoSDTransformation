package main

import (
	"bpmntosd/custom"
	"gitlab.com/teec2/simplified/components/serverextender/client"
	"log"
	"time"
)

// This main program just starts the Server Extender.
// Functionality can be added in the HandleMessage and MyConfiguration classes.
// All commandline parameters are described in the CommandLineOptions Class
func main() {
	log.Println("----------------------------------------->Starting the WebSocket client...")
	startTime := time.Now()
	var wsClient client.WebSocketClient
	go custom.Start(&wsClient, *ClientID, *ClientSecret, *User, *Password, *Addr, *Port)
	time.Sleep(4 * time.Second)

	elapsedTime := time.Since(startTime)
	log.Printf("-----------------------------------------> WebSocket client connected in %s", elapsedTime)

	custom.DefineModels(&wsClient)

	elapsedTime3 := time.Since(startTime)
	log.Printf("-----------------------------------------> Total time %s", elapsedTime3)

}
