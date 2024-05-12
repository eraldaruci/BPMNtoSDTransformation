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
	// parse the flags that are given on the command line or take the default value
	//flag.Parse()
	go custom.Start(&wsClient, *ClientID, *ClientSecret, *User, *Password, *Addr, *Port)
	time.Sleep(4 * time.Second)

	elapsedTime := time.Since(startTime)
	log.Printf("-----------------------------------------> WebSocket client connected in %s", elapsedTime)
	log.Println("----------------------------------------->Starting model-to-model transformation...")
	startTime2 := time.Now()
	custom.DefineModels(&wsClient)
	elapsedTime2 := time.Since(startTime2)
	log.Printf("-----------------------------------------> Model-to-model transformation done in %s", elapsedTime2)
	elapsedTime3 := time.Since(startTime)
	log.Printf("-----------------------------------------> Total time %s", elapsedTime3)

}
