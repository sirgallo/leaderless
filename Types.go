package athn

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/connpool"
import "github.com/sirgallo/athn/liveness"
import "github.com/sirgallo/athn/propose"
import "github.com/sirgallo/athn/request"
import "github.com/sirgallo/athn/system"


type AthnPortOpts struct {
	Liveness int
	Proposal int
	Request int
}

type AthnServiceOpts struct {
	NodeSeed []byte
	Ports AthnPortOpts
	Protocol string
	ConnPoolOpts connpool.ConnectionPoolOpts
}

type Athn [T request.Payload, U request.Result] struct {
	ports AthnPortOpts
	protocol string

	system *system.System
	zLog logger.CustomLog

	livenessService *liveness.LivenessService
	proposeService *propose.ProposeService[T, U]
	requestService *request.RequestService[T, U]
}


const NAME = "Athn"