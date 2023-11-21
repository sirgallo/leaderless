package athn

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/connpool"
import "github.com/sirgallo/athn/liveness"
import "github.com/sirgallo/athn/propose"
import "github.com/sirgallo/athn/request"
import "github.com/sirgallo/athn/system"


func NewAthn[T request.Payload, U request.Result](opts AthnServiceOpts) (*Athn[T, U], error) {
	sys, newSysErr := system.NewSystem(opts.NodeSeed)
	if newSysErr != nil { return nil, newSysErr }

	athn := &Athn[T, U]{
		ports: opts.Ports,
		protocol: opts.Protocol,
		system: sys,
		zLog: *logger.NewCustomLog(NAME),
	}

	lConnPool := connpool.NewConnectionPool(opts.ConnPoolOpts)
	pConnPool := connpool.NewConnectionPool(opts.ConnPoolOpts)

	lOpts := &liveness.LivenessServiceOpts{ 
		Port: opts.Ports.Liveness, 
		ConnectionPool: lConnPool, 
		System: sys,
	}

	pOpts := &propose.ProposeServiceOpts{
		Port: opts.Ports.Proposal,
		ConnectionPool: pConnPool,
		System: sys,
	}

	reqOpts := &request.RequestServiceOpts{ Port: opts.Ports.Request, System: sys }

	athn.livenessService = liveness.NewLivenessService(lOpts)
	athn.proposeService = propose.NewProposeService[T, U](pOpts)
	athn.requestService = request.NewRequestService[T, U](reqOpts)

	return athn, nil
}

func (athn *Athn[T, U]) StartAthn() {
	athn.StartModules()
	athn.StartModulePassThroughs()
	
	select {}
}