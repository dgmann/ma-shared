package lookup

type MemoryStorage struct {
	table map[string]bool
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{table:make(map[string]bool)}
}

func(storage *MemoryStorage) Exists(keys ...string) bool {
	for _, key := range keys {
		_, exists := storage.table[key]
		if !exists {
			return false
		}
	}
	return true
}

func(storage *MemoryStorage) Close() {
	storage.table = make(map[string]bool)
}

func(storage *MemoryStorage) Add(key string) {
	storage.table[key] = true
}
