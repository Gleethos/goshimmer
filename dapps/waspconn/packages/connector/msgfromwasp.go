package connector

import (
	"github.com/iotaledger/goshimmer/dapps/waspconn/packages/chopper"
	"github.com/iotaledger/goshimmer/dapps/waspconn/packages/waspconn"
	"github.com/iotaledger/goshimmer/packages/binary/messagelayer/payload"
)

// process messages received from the Wasp
func (wconn *WaspConnector) processMsgDataFromWasp(data []byte) {
	var msg interface{}
	var err error
	if msg, err = waspconn.DecodeMsg(data, false); err != nil {
		wconn.log.Errorf("DecodeMsg: %v", err)
		return
	}
	switch msgt := msg.(type) {
	case *waspconn.WaspMsgChunk:
		finalMsg, err := chopper.IncomingChunk(msgt.Data, payload.MaxMessageSize-waspconn.ChunkMessageHeaderSize)
		if err != nil {
			wconn.log.Errorf("DecodeMsg: %v", err)
			return
		}
		if finalMsg != nil {
			wconn.processMsgDataFromWasp(finalMsg)
		}

	case *waspconn.WaspPingMsg:
		wconn.log.Debugf("PING %d received", msgt.Id)
		if err := wconn.sendMsgToWasp(msgt); err != nil {
			wconn.log.Errorf("responding to ping: %v", err)
		}

	case *waspconn.WaspToNodeTransactionMsg:
		wconn.postTransaction(msgt.Tx, &msgt.SCAddress, msgt.Leader)

	case *waspconn.WaspToNodeSubscribeMsg:
		for _, addr := range msgt.Addresses {
			wconn.subscribe(&addr)
		}
		go func() {
			for _, addr := range msgt.Addresses {
				wconn.pushBacklogToWasp(&addr)
			}
		}()

	case *waspconn.WaspToNodeGetTransactionMsg:
		wconn.getTransaction(msgt.TxId)

	case *waspconn.WaspToNodeGetOutputsMsg:
		wconn.getAddressBalance(&msgt.Address)

	case *waspconn.WaspToNodeSetIdMsg:
		wconn.SetId(msgt.Waspid)

	default:
		panic("wrong msg type")
	}
}
