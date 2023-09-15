package crawler

import "time"

// Ticker is similar to time.Ticket, except that it accepts an initial delay
type Ticker struct {
	C <-chan time.Time

	stop chan<- bool
}

func NewTicker(initialDelay time.Duration, period time.Duration) *Ticker {
	c := make(chan time.Time)
	stop := make(chan bool)
	go func() {
		defer close(c)
		defer close(stop)

		// wait for the initial delay to elapse using a simple timer
		initialDelayTimer := time.NewTimer(initialDelay)
		defer initialDelayTimer.Stop()
	initialDelayLoop:
		for {
			select {
			case <-stop:
				return
			case nowTime := <-initialDelayTimer.C:
				c <- nowTime
				break initialDelayLoop
			}
		}

		// once the initial delay has elapsed, switch to a regular ticker to tick every period
		periodTicker := time.NewTicker(period)
		defer periodTicker.Stop()
		for {
			select {
			case <-stop:
				return
			case nowTime := <-periodTicker.C:
				c <- nowTime
			}
		}
	}()
	return &Ticker{
		C:    c,
		stop: stop,
	}
}

func (ticker *Ticker) Stop() {
	ticker.stop <- true
}
