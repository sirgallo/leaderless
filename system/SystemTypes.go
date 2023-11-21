package system

import "sync"
// import "unsafe"


type System struct {
	Host string
	NodeId[32]byte
	VersionTag uint64
	// proposalChannel chan *proposal.Proposal
	
	PropagationFactor int

	Neighbors *sync.Map //map[[32]byte] *liveness.NodeInfo
	PreviousNeighbors *sync.Map // map[[32]byte] struct{}

	// updateVersion *unsafe.Pointer
	
	// StateMachine *statemachine.StateMachine

	SystemMutex sync.Mutex
}