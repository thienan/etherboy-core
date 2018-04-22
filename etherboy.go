package main

import (
	"strings"

	proto "github.com/gogo/protobuf/proto"
	contract "github.com/loomnetwork/etherboy-core/contract-helpers"
	"github.com/loomnetwork/etherboy-core/txmsg"
	"github.com/loomnetwork/loom"
	"github.com/loomnetwork/loom/plugin"
	"github.com/pkg/errors"
)

func main() {}

type EtherBoy struct {
	contract.SimpleContract
}

func (e *EtherBoy) Meta() plugin.Meta {
	return plugin.Meta{
		Name:    "etherboycore",
		Version: "0.0.1",
	}
}

func (e *EtherBoy) Init(ctx plugin.Context, req *plugin.Request) error {
	return nil
}

func (e *EtherBoy) CreateAccount(ctx plugin.Context, accTx *txmsg.EtherboyCreateAccountTx) error {
	owner := strings.TrimSpace(accTx.Owner)
	// confirm owner doesnt exist already
	if ctx.Has(e.ownerKey(owner)) {
		return errors.New("Owner already exists")
	}
	state := &txmsg.EtherboyAppState{
		Address: []byte(ctx.Message().Sender.Local),
	}
	statebytes, err := proto.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "Error marshaling state node")
	}
	ctx.Set(e.ownerKey(owner), statebytes)
	return nil
}

func (e *EtherBoy) SaveState(ctx plugin.Context, tx *txmsg.EtherboyStateTx) error {
	owner := strings.TrimSpace(tx.Owner)
	var curState txmsg.EtherboyAppState
	if err := proto.Unmarshal(ctx.Get(e.ownerKey(owner)), &curState); err != nil {
		return err
	}
	if loom.LocalAddress(curState.Address).Compare(ctx.Message().Sender.Local) != 0 {
		return errors.New("Owner unverified")
	}
	curState.Blob = tx.Data
	statebytes, err := proto.Marshal(&curState)
	if err != nil {
		return errors.Wrap(err, "Error marshaling state node")
	}
	ctx.Set(e.ownerKey(owner), statebytes)
	return nil
}

func (s *EtherBoy) ownerKey(owner string) []byte {
	return []byte("owner:" + owner)
}

func NewEtherBoyContract() plugin.Contract {
	e := &EtherBoy{}
	e.SimpleContract.Init(e)
	return e
}

var Contract = NewEtherBoyContract()
