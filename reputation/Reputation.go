package reputation

import "math"


func NewReputionStrategy(opts *ReputationStrategyOpts) *ReputationStrategy {
	return &ReputationStrategy{
		globals: opts.Globals,
		system: opts.System,
	}
}

func (rep *ReputationStrategy) CalculatePreferenceScore(nodeId [32]byte, redundantCopies int) (float64, error) {
	reputationScore, calcErr := rep.CalculateReputation(nodeId)
	if calcErr != nil { return 0, calcErr }

	return WEIGHT_REDUNDANT_PROPOSALS * float64(redundantCopies) + WEIGHT_REPUTATION_SCORE * reputationScore, nil
}

func (rep *ReputationStrategy) CalculateReputation(nodeId [32]byte) (float64, error) {
	version, nodeInfo, readErr := rep.globals.ReadNodeInfoAndVersion(nodeId)
	if readErr != nil { return float64(0), readErr }

	proportionOfInfluence := float64(nodeInfo.SuccessWrites) / float64(version)
	return logBase(proportionOfInfluence, float64(rep.system.PropagationFactor)), nil
}

func logBase(x, base float64) float64 {
	return math.Log(x) / math.Log(base)
}