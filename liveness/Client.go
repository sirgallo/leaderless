package liveness

import "bytes"
import "context"
import "sync"
import "google.golang.org/grpc"

import "github.com/sirgallo/athn/common/serialize"
import "github.com/sirgallo/athn/common/utils"
import "github.com/sirgallo/athn/globals"
import "github.com/sirgallo/athn/proto/liveness"


func (liveSrv *LivenessService) RequestLiveness() error {
	respChans := liveSrv.createResponseChannels()
	var messages []*liveness.LivenessMessage
	var livenessWG sync.WaitGroup

	livenessWG.Add(1)
	go func() {
		defer livenessWG.Done()
		liveSrv.broadcastLivenessMsgs(respChans)
	}()

	livenessWG.Add(1)
	go func() {
		defer livenessWG.Done()

		liveSrv.broadcastLivenessMsgs(respChans)
		for {
			select {
				case <- respChans.BroadcastClose:
					for message := range respChans.Messages {
						messages = append(messages, message)
					}

					return
				case message :=<- respChans.Messages:
					messages = append(messages, message)
			}
		}
	}()

	transform := func(msg *liveness.LivenessMessage) *liveness.NodeInfo { return msg.Sender }
	nodeInfoArr := utils.Map[*liveness.LivenessMessage, *liveness.NodeInfo](messages, transform)

	_, updateErr := liveSrv.system.UpdateNeighbors(nodeInfoArr)
	if updateErr != nil { return updateErr }

	return nil
}

func (liveSrv *LivenessService) broadcastLivenessMsgs(respChans LivenessResponseChannels) error {
	defer close(respChans.BroadcastClose)
	defer close(respChans.Messages)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	version, readErr := liveSrv.system.Globals.GetVersion()
	if readErr != nil { return readErr }

	sVersion := serialize.SerializeBigInt(version, globals.GLOBAL_V_BYTE_LENGTH)
	var broadcastWG sync.WaitGroup

	message := &liveness.LivenessMessage{
		VersionTag: sVersion,
		Sender: &liveness.NodeInfo{
			Host: liveSrv.system.Host,
			NodeId: liveSrv.system.NodeId[:],
			OK: true,
		},
	}

	liveSrv.system.Neighbors.Range(func(key any, value any) bool {
		neighbor := value.(*liveness.NodeInfo)
		
		broadcastWG.Add(1)

		go func(neighbor *liveness.NodeInfo) {
			defer broadcastWG.Done()
			
			conn, connErr := liveSrv.connPool.GetConnection(neighbor.Host, liveSrv.Port)
			if connErr != nil {
				liveSrv.zLog.Error("Failed to connect to", neighbor.Host + liveSrv.Port, ":", connErr.Error())
				return
			}

			select {
				case <- ctx.Done():
					liveSrv.connPool.PutConnection(neighbor.Host, conn)
					return
				default:
					res, _ := liveSrv.clientLivenessRPC(conn, sVersion, message, neighbor)
					respChans.Messages <- res

					if bytes.Compare(res.VersionTag, sVersion) == 1 {
						liveSrv.zLog.Debug("higher version found on response")
					}

					liveSrv.connPool.PutConnection(neighbor.Host, conn)
			}
		}(neighbor)

		return true
	})

	broadcastWG.Wait()
	return nil
}

func (liveSrv *LivenessService) clientLivenessRPC(
	conn *grpc.ClientConn,
	sVersion []byte,
	message *liveness.LivenessMessage,
	neighbor *liveness.NodeInfo,
) (*liveness.LivenessMessage, error) {
	client := liveness.NewLivenessClient(conn)

	livenessRPC := func() (*liveness.LivenessMessage, error) {
		ctx, cancel := context.WithTimeout(context.Background(), LIVENESS_RPC_TIMEOUT)
		defer cancel()

		res, err := client.LivenessRPC(ctx, message)
		if err != nil {
			liveSrv.zLog.Error("exp backoff attempt err:", err.Error())
			return utils.GetZero[*liveness.LivenessMessage](), err 
		}

		return res, nil
	}

	maxRetries := 3
	expOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMilliseconds: 10 }
	expBackoff := utils.NewExponentialBackoffStrat[*liveness.LivenessMessage](expOpts)

	res, err := expBackoff.PerformBackoff(livenessRPC)
	if err != nil {
		liveSrv.zLog.Warn("system", neighbor.Host, "unreachable, setting status to dead")
		liveSrv.connPool.CloseConnections(neighbor.Host)

		res = &liveness.LivenessMessage{
			VersionTag: sVersion,
			Sender: &liveness.NodeInfo{ Host: neighbor.Host, NodeId: neighbor.NodeId, OK: false },
		}
	}

	return res, nil
}

func (liveSrv *LivenessService) createResponseChannels() LivenessResponseChannels {
	broadcastClose := make(chan struct{})
	messages := make(chan *liveness.LivenessMessage, liveSrv.system.PropagationFactor)

	return LivenessResponseChannels{
		BroadcastClose: broadcastClose,
		Messages: messages,
	}
}