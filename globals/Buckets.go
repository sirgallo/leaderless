package globals

import "math/big"


const GLOBALS_FILENAME = "globals.db"
const GLOBALS_SUB_DIRECTORY = "globals"

const GLOBAL_BUCKET = "global"

const GLOBAL_VERSION_BUCKET = "version"
const GLOBAL_VERSION_KEY = "current"

type GlobalVersionValue = big.Int

const GLOBAL_NODEINFO_BUCKET = "nodeinfo"

type GlobalNodeInfoKey = [32]byte
type GlobalNodeInfoValue struct {
	SuccessWrites uint64
	OK bool
}

const (
	SuccessWritesOffset = 0
	NodeOKOffset = 8
)

const GLOBAL_V_BYTE_LENGTH = 128