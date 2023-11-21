package athn

import "net"


func (athn *Athn[T, U]) StartModules() {
	lListener, leErr := net.Listen(athn.protocol, athn.livenessService.Port)
	if leErr != nil { athn.zLog.Error("Failed to listen: %v", leErr.Error()) }

	pListener, rlErr := net.Listen(athn.protocol, athn.proposeService.Port)
	if rlErr != nil { athn.zLog.Error("Failed to listen: %v", rlErr.Error()) }

	go athn.livenessService.StartLivenessService(&lListener)
	go athn.proposeService.StartProposeService(&pListener)
}

func (athn *Athn[T, U]) StartModulePassThroughs() {
	go func() {
		for request := range athn.requestService.RequestBuffer {
			athn.proposeService.ClientReqBuffer <- request
		}
	}()

	go func() {
		for response := range athn.proposeService.ClientRespBuffer {
			athn.requestService.ResponseBuffer <- response
		}
	}()
}