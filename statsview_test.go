package statsview

import (
	"testing"
	"time"
)

func TestStatsViewMgr(t *testing.T) {
	timeout := time.After(time.Minute)
	done := make(chan bool)

	go func() {
		mgr := New()
		go mgr.Start()
		time.Sleep(10 * time.Second)
		mgr.Stop()

		time.Sleep(2 * time.Second)

		mgr = New()
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
