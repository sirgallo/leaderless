package request

import "net/http"
import "sync"
import "time"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/system"


type Payload = comparable
type Result = comparable

type RequestServiceOpts struct {
	Port int
	System *system.System
}

type RequestService[T Payload, U Result] struct {
	mux *http.ServeMux
	
	Port string
	system *system.System
	zLog logger.CustomLog
	
	clientMappedResponseChannels sync.Map
	RequestBuffer chan *ClientRequest[T]
	ResponseBuffer chan *ClientResponse[U]
}

type ClientRequest[T Payload] struct {
	RequestID [32]byte `json:"-"`
	Payload T `json:"payload"`
}

type ClientResponse[T Result] struct {
	RequestID [32]byte `json:"-"`
	Result T `json:"result"`
	Success bool `json:"success"`
	ErrorMsg *error `json:"error_msg"`
}

const NAME = "Request Service"
const CommandRoute = "/command"
const RequestChannelSize = 1000000
const ResponseChannelSize = 1000000
const HTTPTimeout = 2 * time.Second