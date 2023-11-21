package system

import "errors"
import "os"
import "runtime"
import "sync/atomic"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/utils"
import "github.com/sirgallo/athn/proto/liveness"


//=========================================== System


const NAME = "System"
var Log = logger.NewCustomLog(NAME)

func NewSystem(seed []byte) (*System, error) {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil { return nil, hostnameErr }

	nodeId, hashErr := utils.GenerateSHA256HashWithSeed(seed)
	if hashErr != nil { return nil, hashErr }

	return &System{
		Host: hostname,
		NodeId: nodeId,
	}, nil
}

func (sys *System) SetHost(host string) bool {
	sys.SystemMutex.Lock()
	defer sys.SystemMutex.Unlock()

	sys.Host = host
	return true
}

func (sys *System) SetNodeId(nodeId [32]byte) bool {
	sys.SystemMutex.Lock()
	defer sys.SystemMutex.Unlock()

	sys.NodeId = nodeId
	return true
}

func (sys *System) UpdateVersionTag(updatedVTag uint64) bool {
	for updatedVTag == uint64(sys.VersionTag) + 1 &&
	 ! atomic.CompareAndSwapUint64(&sys.VersionTag, sys.VersionTag, updatedVTag) {
		runtime.Gosched()
	}

	return true
}

func (sys *System) UpdateNeighbors(newNeighbors []*liveness.NodeInfo) (bool, error) {
	sys.SystemMutex.Lock()
	defer sys.SystemMutex.Unlock()

	sys.PreviousNeighbors.Range(func(key any, value any) bool {
		sys.PreviousNeighbors.Delete(key.([32]byte))
		return true
	})

	sys.Neighbors.Range(func(key any, value any) bool {
		sys.Neighbors.Delete(key.([32]byte))
		sys.PreviousNeighbors.Store(key.([32]byte), struct{}{})
		
		return true
	})

	for _, newNeighbor := range newNeighbors {
		nodeId, getNodeIdErr := func() ([32]byte, error) {
			if len(newNeighbor.NodeId) > 32 { 
				return utils.GetZero[[32]byte](), errors.New("neighbor nodeId incorrect length, should be 32 bytes") 
			}

			var arr [32]byte
			copy(arr[:], newNeighbor.NodeId)
			return arr, nil
		}()

		if getNodeIdErr != nil { return false, getNodeIdErr }

		sys.Neighbors.Store(nodeId, newNeighbor)
	}
	
	return true, nil
}