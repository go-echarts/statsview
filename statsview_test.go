package statsview

import (
	"testing"
	"time"
)

func TestStatsViewMgr(t *testing.T) {
	mgr := New()

	timeout := time.After(12 * time.Second)
	done := make(chan bool)
	go func() {
		go mgr.Start()
		time.Sleep(10 * time.Second)
		mgr.Stop()
		
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}
