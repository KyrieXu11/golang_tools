package constant

const (
	Deleted = -1
	InUse   = 0

	CommandList = iota
	CommandUse
	CommandSet
	CommandGet
	CommandDelete

	InvalidId = -1
)
