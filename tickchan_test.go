package tickchan

import (
	"testing"
	"time"
)

func TestChainTicker_Add(t *testing.T) {
	ticker := NewChanTicker(time.Millisecond * 1)
	ch := make(chan time.Time, 1)
	ticker.Add(ch)
	t0 := time.Now()
	for i := 0; i < 100; i++ {
		<-ch
	}
	t.Logf("time used: %s, expect %s", time.Now().Sub(t0), time.Millisecond*100)
}

func TestChainTicker_Now(t *testing.T) {
	now := time.Now()
	ticker := NewChanTicker(time.Millisecond)
	if ticker.LastTickTime().Before(now) {
		t.Fatalf("ticker LastTickTime is wrong")
	}

	for i := 0; i < 1000; i++ {
		if time.Now().Sub(ticker.LastTickTime()) > time.Millisecond*10 {
			t.Fatalf("ticker.LastTickTime() is too old")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestChainTicker_Remove(t *testing.T) {
	ticker := NewChanTicker(time.Millisecond)
	ch := make(chan time.Time, 1)
	ticker.Add(ch)
	for i := 0; i < 10; i++ {
		<-ch
	}

	ticker.Remove(ch)
	select {
	case <-ch:
		t.Fatalf("should not receive from channel")
	case <-time.After(time.Millisecond * 100):
		break
	}
}
