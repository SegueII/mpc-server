package main

import (
	sdk "github.com/irisnet/irishub-sdk-go"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/store"

	"github.com/segueII/mpc-server/aby"
)

var (
	nodeURI  = "tcp://localhost:26657"
	grpcAddr = "localhost:9090"
	chainID  = "testing"
)

const (
	// addr     = "iaa179jc96sqxfyegr4vtwgu3gr32q42xn8mkefver"
	keyName  = "key"
	password = "12345678"
	mnemonic = "small pretty lock logic loud beef please boring space picnic essence also opera come roast pepper pumpkin vivid topple asset upon dismiss debris awful"
)

func main() {
	client := initClient()

	setKey(client, keyName, password, mnemonic)

	subscribeServiceRequest(client, "test", "iaa179jc96sqxfyegr4vtwgu3gr32q42xn8mkefver")
}

func initClient() sdk.IRISHUBClient {
	options := []types.Option{
		types.KeyDAOOption(store.NewMemory(nil)),
		types.TimeoutOption(10),
	}

	cfg, err := types.NewClientConfig(nodeURI, grpcAddr, chainID, options...)
	if err != nil {
		panic(err)
	}

	return sdk.NewIRISHUBClient(cfg)
}

func setKey(client sdk.IRISHUBClient, name string, password string, mnemonic string) {
	_, _ = client.Key.Recover(name, password, mnemonic)
}

func subscribeServiceRequest(client sdk.IRISHUBClient, serviceName string, provider string) {
	doneChan := make(chan struct{})
	ch := make(chan string)
	sub, _ := client.Service.SubscribeServiceRequest(
		serviceName,
		func(reqCtxID, reqID, input string) (output string, result string) {
			abyClient := aby.NewABY()
			out, err := abyClient.Server("env")
			if err != nil {
				panic(err)
			}
			println(string(out))

			ch <- reqID
			output = `{"header":{},"body":{"test":"xxx"}}`
			result = `{"code":200,"message":""}`
			return output, result
		},
		types.BaseTx{
			From:     keyName,
			Gas:      200000,
			Memo:     "",
			Mode:     types.Sync,
			Password: password,
		},
	)

	for {
		select {
		case reqID := <-ch:
			println(reqID)
		case <-doneChan:
			_ = client.Unsubscribe(sub)
			println("done")
		}
	}
}
