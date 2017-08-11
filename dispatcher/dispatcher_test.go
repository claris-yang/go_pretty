package dispatcher

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	dispatcher := NewDispatcher(2000)
	dispatcher.Run()

	//runtime.GOMAXPROCS(8)
	for i := 0; i < 3000; i++ {
		job := NewJob(i)
		dispatcher.AddJob(job)
		//go job.Handle()
	}

	go func() {
		for {
			select {
			case <-time.After(time.Minute * 5):
				fmt.Println(" dispatcher. \n ", dispatcher.Count)
				wg.Done()
				break
			}
		}
	}()

	wg.Wait()
	fmt.Printf("well done!")

}
