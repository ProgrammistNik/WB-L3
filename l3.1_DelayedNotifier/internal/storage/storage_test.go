package storage

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.1/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestStorage_SetAndGet(t *testing.T) {
	st := New()

	notif := model.Notification{
		ID:        "test-id",
		Message:   "test message",
		SendAt:    time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		Status:    "scheduled",
	}

	st.Set(notif)

	got, ok := st.Get("test-id")
	assert.True(t, ok)
	assert.Equal(t, notif, got)
}

func TestStorage_GetNotFound(t *testing.T) {
	st := New()

	_, ok := st.Get("missing-id")
	assert.False(t, ok)
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	st := New()
	wg := sync.WaitGroup{}

	n := 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := "id-" + strconv.Itoa(i)
			notif := model.Notification{ID: id}
			st.Set(notif)
			got, ok := st.Get(id)
			assert.True(t, ok)
			assert.Equal(t, id, got.ID)
		}(i)
	}

	wg.Wait()
}
