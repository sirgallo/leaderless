package athn

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/connpool"
import "github.com/sirgallo/athn/liveness"
import "github.com/sirgallo/athn/propose"
import "github.com/sirgallo/athn/request"


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

type Athn struct {
	ports AthnPortOpts
	protocol string
	zLog logger.CustomLog

	livenessService *liveness.LivenessService
	proposeService *propose.ProposeService
	requestService *request.RequestService
}


const NAME = "Athn"