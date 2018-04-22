package main

import (
	"fmt"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/loomnetwork/etherboy-core/txmsg"
	"github.com/loomnetwork/loom/plugin"
)

func TestCreateAccount(t *testing.T) {
	e := &EtherBoy{}
	tx := &txmsg.EtherboyCreateAccountTx{
		Version: 0,
		Owner:   "aditya",
		Data:    []byte("dummy"),
	}

	any, err := ptypes.MarshalAny(tx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	msg := &txmsg.SimpleContractMethod{
		Version: 0,
		Method:  fmt.Sprintf("%s.CreateAccount", e.Meta().Name),
		Data:    any,
	}

	msgBytes, err := proto.Marshal(msg)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	ctx := CreateFakeContext()
	req := &plugin.Request{}

	fmt.Printf("Data: msgBytes-%v\n", msgBytes)
	err = e.Init(ctx, req)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	req.Body = msgBytes
	resp, err := e.Call(ctx, req)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	fmt.Printf("resp -%v\n", resp)
}