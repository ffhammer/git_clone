package treediff

type ChangeAction string

const (
	Add    ChangeAction = "added"
	Modify ChangeAction = "modified"
	Delete ChangeAction = "deleted"
	// Stash    ChangeAction = "stash"
	// Unmerged ChangeAction = "unmerged"
)

type ChangeEntry struct {
	RelPath    string       `json:"relpath"`
	NewHash    string       `json:"filehash"`
	OldHash    string       `json:"oldHash"`
	EditedTime int64        `json:"editTime"`
	Action     ChangeAction `json:"actiion"`
}
type ChangeCollector[T any] interface {
	// Add a change to the collector.
	Add(change ChangeEntry)
	// Return the final container.
}
type ChangeMap map[string]ChangeEntry

func (cm ChangeMap) Add(change ChangeEntry) {
	cm[change.RelPath] = change
}

type ChangeList []ChangeEntry

func (cl *ChangeList) Add(change ChangeEntry) {
	*cl = append(*cl, change)
}
