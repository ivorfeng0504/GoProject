package dbsync

type SyncResult struct {
	TableID string
	State   bool
	Message string
}
