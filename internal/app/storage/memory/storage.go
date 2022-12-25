package memory

import (
	"log"
	"sync"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type Memory struct {
	mu         sync.RWMutex
	storage    map[string]*storage.LocalShortenData
	repository storage.Repository
}

func New() (*Memory, error) {
	return &Memory{
		storage: map[string]*storage.LocalShortenData{},
	}, nil
}

func (m *Memory) Init() error {
	m.repository.InputCh = make(chan storage.DataForWorker, 100)
	return nil
}

func (m *Memory) SaveURL(id string, data *storage.LocalShortenData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[id] = data
	return nil
}

func (m *Memory) GetURL(id string) (*storage.LocalShortenData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.storage[id]
	if !ok {
		return &storage.LocalShortenData{}, storage.ErrNotFound
	}

	return v, nil
}

func (m *Memory) GetByField(field, val string) (*storage.LocalShortenData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch field {
	case "id":
		v, ok := m.storage[val]
		if !ok {
			return &storage.LocalShortenData{}, storage.ErrNotFound
		}
		return v, nil
	case "url":
		for _, v := range m.storage {
			if v.URL == val {
				return v, nil
			}
		}
		return &storage.LocalShortenData{}, storage.ErrNotFound
	default:
		return &storage.LocalShortenData{}, storage.ErrNotFound
	}
}

func (m *Memory) GetAllURLsForSignID(signID uint32) ([]*storage.LocalShortenData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var (
		urls []*storage.LocalShortenData
	)
	for _, item := range m.storage {
		if item.SignID == signID {
			urls = append(urls, item)
		}
	}

	if len(urls) < 1 {
		return nil, storage.ErrNotFound
	}

	return urls, nil
}

func (m *Memory) DeleteURL(id string, val bool, signID uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.storage[id]
	if !ok {
		return storage.ErrNotFound
	}

	if v.SignID != signID {
		return storage.ErrNotAccepted
	}
	m.storage[id].IsDeleted = val
	return nil
}

func (m *Memory) RunWorkers(count int) {
	for i := 0; i < count; i++ {
		m.repository.WG.Add(1)
		go func() {
			for {
				data, ok := <-m.repository.InputCh
				if !ok {
					log.Printf("Канал закрылся, завершаем работу")
					m.repository.WG.Done()
					return
				}
				log.Printf("Отправляем данные в БД!")
				err := m.DeleteURL(data.ID, true, data.SignID)
				if err != nil {
					log.Printf("Проблема в бд, %v", err)
				}
			}
		}()
	}
}

func (m *Memory) AddJob(urlIDs []string, signID uint32) {
	go func() {
		for _, urlID := range urlIDs {
			m.repository.InputCh <- storage.DataForWorker{
				ID:     urlID,
				SignID: signID,
			}
		}
	}()
}

func (m *Memory) Stop() {
	close(m.repository.InputCh)
}

func (m *Memory) Wait() {
	m.repository.WG.Wait()
}
