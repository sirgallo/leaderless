package globals


const GLOBALS_FILENAME = "globals.db"
const GLOBALS_SUB_DIRECTORY = "globals"

const GLOBAL_BUCKET = "global"

const GLOBAL_VERSION_BUCKET = "version"
const GLOBAL_VERSION_KEY = "current"

type GlobalVersionValue = uint64

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