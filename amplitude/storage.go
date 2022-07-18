package amplitude

type Storage interface {
	push()
	pull()
	pullAll()
}

type StorageProvider interface {
	GetStorage() Storage
}

type InMemoryStorage struct {
	totalEvents   int
	configuration Config
	workers       Worker
}

func (i *InMemoryStorage) setup(configuration Config, workers Worker) {
}

func (i InMemoryStorage) push(event Event, delay int) {

}

func (i InMemoryStorage) pull(batchSize int) {

}

func (i InMemoryStorage) pullAll() {

}

func (i InMemoryStorage) insertEvent(totalDelay int, event Event) {

}

func (i InMemoryStorage) getRetryDelay(retry int) {

}

type InMemoryStorageProvider struct {
}

func (i InMemoryStorageProvider) getStorage() {

}
