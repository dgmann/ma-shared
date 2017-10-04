package lookup

type StorageClient interface {
	Exists(key ...string) bool
	Add(key string)
	Close()
}
