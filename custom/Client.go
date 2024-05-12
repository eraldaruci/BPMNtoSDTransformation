package custom

import "gitlab.com/teec2/simplified/components/serverextender/client"

func Start(wsClient *client.WebSocketClient, clientId string, clientSecret string, userName string, password string, address string, port int) {
	client.RunServerExtender(wsClient, client.Credentials{
		User:         userName,
		Password:     password,
		ClientSecret: clientSecret,
		ClientId:     clientId,
	}, client.ModuleConfig{
		ExtServer:   ExtMyServer,
		ExtServerId: ExtMyServerId,
		ExtModule:   ExtMyModule,
	}, address, port, HandleMessage, Configuration)
}
