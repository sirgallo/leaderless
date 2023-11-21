package athn

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/connpool"
import "github.com/sirgallo/athn/globals"
import "github.com/sirgallo/athn/liveness"
import "github.com/sirgallo/athn/propose"
import "github.com/sirgallo/athn/request"
import "github.com/sirgallo/athn/system"
import "github.com/sirgallo/athn/state"


func NewAthn(opts AthnServiceOpts) (*Athn, error) {
	globals, openGlobalsErr := globals.NewGlobals()
	if openGlobalsErr != nil { return nil, openGlobalsErr }

	state, openStateErr := state.NewState()
	if openStateErr != nil { return nil, openGlobalsErr }
	
	sysOpts := &system.SystemOpts{
		Seed: opts.NodeSeed,
		Globals: globals,
		State: state,
	}

	sys, newSysErr := system.NewSystem(sysOpts)
	if newSysErr != nil { return nil, newSysErr }

	athn := &Athn{
		ports: opts.Ports,
		protocol: opts.Protocol,
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
	athn.proposeService = propose.NewProposeService(pOpts)
	athn.requestService = request.NewRequestService(reqOpts)

	return athn, nil
}

func (athn *Athn) StartAthn() {
	athn.StartModules()
	athn.StartModulePassThroughs()
	
	select {}
}