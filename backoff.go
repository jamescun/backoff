package backoff

import (
	"math"
	"math/rand"
	"net"
	"time"
)

// random source for jitter (if enabled), no need to be cryptographically random
var backoffRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type Backoff struct {
	// minimum backoff time
	Minimum time.Duration

	// maximum number of failures to tolerate
	Maximum int

	// add randomness to backoff duration
	// http://www.awsarchitectureblog.com/2015/03/backoff.html
	Jitter bool
}

// create a new exponential Backoff assignment with minimum time between
// failures and maximum total failures.
func New(min time.Duration, max int) Backoff {
	return Backoff{
		Minimum: min,
		Maximum: max,
	}
}

func (b Backoff) duration(failures int) time.Duration {
	d := int64(b.Minimum) * int64(math.Pow(2, float64(failures)))

	if b.Jitter {
		d = (d / 2) + backoffRand.Int63n((d / 2))
	}

	return time.Duration(d)
}

// execute function fn with exponential backoff mechanism.
// will return last error returned by fn, or nil on first successful attempt.
func (b Backoff) Do(fn func() error) (err error) {
	var failures int

	for {
		err = fn()
		if err == nil {
			break
		}

		failures++
		if failures >= b.Maximum {
			break
		}
		time.Sleep(b.duration(failures))
	}

	return
}

// implements dialer interface with exponential backoff mechanism.
func (b Backoff) Dial(network, address string) (conn net.Conn, err error) {
	b.Do(func() error {
		conn, err = net.Dial(network, address)
		return err
	})

	return
}
