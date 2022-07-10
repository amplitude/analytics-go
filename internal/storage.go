package internal

import (
	"github.com/amplitude/Amplitude-Go/pkg/amplitude"
)

type Storage interface {
	push()
	pull()
	pullAll()
}

type StorageProvider interface {
	GetStorage() Storage
}

type InMemoryStorage struct {
	totalEvents int
	//bufferData
	//readyQueue
	//bufferLockCv
	configuration amplitude.Config
	workers       Worker
}

func (i *InMemoryStorage) setup(configuration amplitude.Config, workers Worker) {
}

func (i InMemoryStorage) push(event amplitude.BaseEvent, delay int) {

}

func (i InMemoryStorage) pull(batchSize int) {

}

func (i InMemoryStorage) pullAll() {

}

func (i InMemoryStorage) insertEvent(totalDelay int, event amplitude.BaseEvent) {

}

func (i InMemoryStorage) getRetryDelay(retry int) {

}

type InMemoryStorageProvider struct {
}

func (i InMemoryStorageProvider) getStorage() {

}
