package propose

import "time"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/connpool"
import "github.com/sirgallo/athn/proto/proposal"
import "github.com/sirgallo/athn/request"
import "github.com/sirgallo/athn/system"


type ProposeServiceOpts struct {
	Port int
	ConnectionPool *connpool.ConnectionPool
	System *system.System
}

type ProposeService struct {
	proposal.UnimplementedProposalServer
	
	Port string
	connPool *connpool.ConnectionPool
	system *system.System
	zLog logger.CustomLog

	ClientReqBuffer chan *request.ClientRequest
	ClientRespBuffer chan *request.ClientResponse
}

type ProposeResponseChannels struct {
	BroadcastClose chan struct{}
}

const NAME = "Proposal Service"

const (
	CLIENT_REQ_BUFFER = 100000
	CLIENT_RESP_BUFFER = CLIENT_REQ_BUFFER
	PROPOSE_RPC_TIMEOUT = 500 * time.Millisecond
)