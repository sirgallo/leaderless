package request

import "net/http"
import "sync"
import "time"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/state"
import "github.com/sirgallo/athn/system"

type Payload = comparable
type Result = comparable

type RequestServiceOpts struct {
	Port int
	System *system.System
}

type RequestService struct {
	mux *http.ServeMux
	
	Port string
	system *system.System
	zLog logger.CustomLog
	
	clientMappedResponseChannels sync.Map
	RequestBuffer chan *ClientRequest
	ResponseBuffer chan *ClientResponse
}

type ClientRequest struct {
	RequestID [32]byte `json:"-"`
	Payload state.StatePayload `json:"payload"`
}

type ClientResponse struct {
	RequestID [32]byte `json:"-"`
	Result state.KeyValuePair `json:"result"`
	Success bool `json:"success"`
	ErrorMsg *error `json:"error_msg"`
}

const NAME = "Request Service"
const CommandRoute = "/command"
const RequestChannelSize = 1000000
const ResponseChannelSize = 1000000
const HTTPTimeout = 2 * time.Second