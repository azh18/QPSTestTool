package task

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/prometheus/common/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type FabricTask struct {
	sdk *fabsdk.FabricSDK

}

func NewFabricTask() (*FabricTask, error) {
	conf := config.FromFile("config_e2e.yaml")
	sdk, err := fabsdk.New(conf)
	if err != nil{
		log.Fatal("create msp failed. err=%+=v", err)
		return nil, err
	}
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("org1.lychee.com"))
	if err != nil{
		log.Fatal("create msp failed. err=%+=v", err)
		return nil, err
	}
	adminIdentity, err := mspClient.GetSigningIdentity("Admin")
	if err != nil {
		log.Fatalf("get admin identify fail: %s\n", err.Error())
	} else {
		fmt.Println("AdminIdentify is found:")
		fmt.Println(adminIdentity)
	}
	channelProvider := sdk.ChannelContext("sunshine",
		fabsdk.WithUser("Admin"),
		fabsdk.WithOrg("member1.example.com"))

	channelClient, err := channel.New(channelProvider)
	if err != nil {
		log.Fatalf("create channel client fail: %s\n", err.Error())
	}

	var args [][]byte
	args = append(args, []byte("a"))

	request := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "query",
		Args:        args,
	}
	response, err := channelClient.Query(request)
	if err != nil {
		log.Fatal("query fail: ", err.Error())
	} else {
		fmt.Printf("response is %s\n", response.Payload)
	}
	return &FabricTask{
		nil,
	}, nil
}

func (f *FabricTask) Do() error {
	return nil
}

func (f *FabricTask) Stop() error{
	f.sdk.Close()
	return nil
}