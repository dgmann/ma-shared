package lookup

type MemoryStorage struct {
	table map[string]bool
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{table:make(map[string]bool)}
}

func(storage *MemoryStorage) Exists(key string) bool {
	_, exists := storage.table[key]
	return exists
}

func(storage *MemoryStorage) Close() {
	storage.table = make(map[string]bool)
}

func(storage *MemoryStorage) Add(key string) {
	storage.table[key] = true
}
