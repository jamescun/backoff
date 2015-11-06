package backoff

import (
	"math"
	"math/rand"
	"net"
	"time"
)

type Backoff struct {
	// minimum backoff time
	Minimum time.Duration

	// maximum number of failures to tolerate
	Maximum int

	// add randomness to backoff duration
	// http://www.awsarchitectureblog.com/2015/03/backoff.html
	Jitter bool

	failures float64
}

// allocate new exponential backoff
func New(min time.Duration, max int) *Backoff {
	return &Backoff{
		Minimum: min,
		Maximum: max,
	}
}

// calculate next duration of exponential backoff
func (b *Backoff) Duration() time.Duration {
	d := b.duration(float64(b.Minimum), b.failures)

	b.failures++
	return d
}

func (b *Backoff) duration(min, failures float64) time.Duration {
	d := min * math.Pow(2, failures)

	if b.Jitter {
		d = rand.Float64()*(d-min) + min
	}

	return time.Duration(d)
}

// implements dialer interface
// connects to server with exponential backoff
// unlike Duration(), Dial() can be used concurrently between many operations
func (b *Backoff) Dial(network, address string) (conn net.Conn, err error) {
	var failures float64
	for {
		conn, err = net.Dial(network, address)
		if err == nil {
			return conn, nil
		}

		failures++
		time.Sleep(b.duration(float64(b.Minimum), failures))
	}

	return nil, err
}

// reset failure count
func (b *Backoff) Reset() {
	b.failures = 0
}
