package tickchan

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type ChainTicker struct {
	chansMu sync.Mutex
	chans   map[chan time.Time]struct{}
	now     unsafe.Pointer
}

func (t *ChainTicker) Add(ch chan time.Time) {
	t.chansMu.Lock()
	t.chans[ch] = struct{}{}
	t.chansMu.Unlock()
}

func (t *ChainTicker) Remove(ch chan time.Time) {
	t.chansMu.Lock()
	delete(t.chans, ch)
	t.chansMu.Unlock()

	for {
		select {
		case <-ch: // keep reading from channel until the channel is empty
		default:
			return
		}
	}
}

func (t *ChainTicker) LastTickTime() time.Time {
	return *(*time.Time)(atomic.LoadPointer(&t.now))
}

func NewChanTicker(tickInterval time.Duration) *ChainTicker {
	createTime := time.Now()
	t := &ChainTicker{
		chans: make(map[chan time.Time]struct{}),
	}
	atomic.StorePointer(&t.now, unsafe.Pointer(&createTime))
	tt := time.NewTicker(tickInterval)

	go func() {
		for now := range tt.C {
			tickTime := now
			atomic.StorePointer(&t.now, unsafe.Pointer(&tickTime))
			t.chansMu.Lock()
			for ch := range t.chans {
				select {
				case ch <- now:
				default:
				}
			}
			t.chansMu.Unlock()
		}
	}()

	return t
}
