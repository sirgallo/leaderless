package reputation

import "github.com/sirgallo/athn/globals"
import "github.com/sirgallo/athn/system"


type ReputationStrategyOpts struct {
	Globals *globals.Globals
	System *system.System
}

type ReputationStrategy struct {
	globals *globals.Globals
	system *system.System
}

const WEIGHT_REDUNDANT_PROPOSALS = float64(0.5)
const WEIGHT_REPUTATION_SCORE = float64(0.5)