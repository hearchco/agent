package search

import (
	"sync"
)

// Waits on either c.Wait() or wg.Wait() to do final.Done().
func waitForSuccessOrFinish(c *sync.Cond, wg *sync.WaitGroup, final *sync.WaitGroup) {
	defer final.Done()
	d := sync.Cond{L: &sync.Mutex{}}

	// Wait for signal from any successful worker.
	go func() {
		c.L.Lock()
		c.Wait()
		c.L.Unlock()

		d.L.Lock()
		d.Signal()
		d.L.Unlock()
	}()

	// Wait for all workers to finish (even if it's unsuccessful).
	go func() {
		wg.Wait()

		d.L.Lock()
		d.Signal()
		d.L.Unlock()
	}()

	// Whichever of the above finishes first, signal the final wait group.
	d.L.Lock()
	d.Wait()
	d.L.Unlock()
}
