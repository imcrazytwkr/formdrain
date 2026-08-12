package cursors

type cursorSort byte

const (
	cursorSortDefault cursorSort = iota
	cursorSortID
	cursorSortHostname
	cursorSortSiteID
	cursorSortCreatedAt
	cursorSortField
)
