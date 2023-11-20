package system

import "sync"
// import "unsafe"

import "github.com/sirgallo/athn/proto/liveness"
// import "github.com/sirgallo/athn/proto/proposal"


type System struct {
	host [32]byte
	versionTag uint64
	// proposalChannel chan *proposal.Proposal
	neighbors map[[32]byte] *liveness.NodeInfo


	// updateVersion *unsafe.Pointer
	
	// StateMachine *statemachine.StateMachine

	SystemMutex sync.Mutex
}