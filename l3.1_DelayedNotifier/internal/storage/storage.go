package storage

import (
	"sync"

	"github.com/ProgrammistNik/WB-L3/l3.1/internal/model"
)

type Storage struct {
	mu            sync.RWMutex
	notifications map[string]model.Notification
}

func New() *Storage {
	return &Storage{
		notifications: make(map[string]model.Notification),
	}
}

func (st *Storage) Set(notif model.Notification) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.notifications[notif.ID] = notif
}

func (st *Storage) Get(id string) (model.Notification, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	notification, ok := st.notifications[id]

	return notification, ok
}
 