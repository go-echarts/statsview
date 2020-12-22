package statsview

import (
	"fmt"
	"testing"
	"time"
)

func TestStatsViewMgr(t *testing.T) {

	timeout := time.After(30 * time.Second)
	done := make(chan bool)
	go func() {

		mgr := New()
		go func() {
			 mgr.Start()

		}()
		time.Sleep(10 * time.Second)
		mgr.Stop()


		time.Sleep(5 * time.Second)
		mgr = New()
		err := mgr.Start()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("123")
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}
