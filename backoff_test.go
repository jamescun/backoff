package backoff

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/pusher/buddha/tcptest"
	"github.com/stretchr/testify/assert"
)

func init() {
	// make randomness source deterministic for testing
	backoffRand = rand.New(rand.NewSource(1))
}

func TestNew(t *testing.T) {
	b := New(100*time.Millisecond, 5)

	assert.Equal(t, b.Minimum, 100*time.Millisecond)
	assert.Equal(t, b.Maximum, 5)
}

func TestBackoffDuration(t *testing.T) {
	b := Backoff{Minimum: 100 * time.Millisecond, Maximum: 5, Jitter: false}

	assert.Equal(t, b.duration(0), 100*time.Millisecond)
	assert.Equal(t, b.duration(1), 200*time.Millisecond)
	assert.Equal(t, b.duration(2), 400*time.Millisecond)
	assert.Equal(t, b.duration(3), 800*time.Millisecond)
	assert.Equal(t, b.duration(4), 1600*time.Millisecond)
}

func BenchmarkBackoffDuration(b *testing.B) {
	bk := Backoff{Minimum: 100 * time.Millisecond, Maximum: 5, Jitter: false}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bk.duration(1)
	}
}

func TestBackoffDurationJitter(t *testing.T) {
	b := Backoff{Minimum: 100 * time.Millisecond, Maximum: 5, Jitter: true}

	// randomness has been made deterministic in init()
	assert.Equal(t, b.duration(0), time.Duration(97779410))
	assert.Equal(t, b.duration(1), time.Duration(182153551))
	assert.Equal(t, b.duration(2), time.Duration(266145821))
	assert.Equal(t, b.duration(3), time.Duration(635010051))
	assert.Equal(t, b.duration(4), time.Duration(1087113937))
}

func BenchmarkBackoffDurationJitter(b *testing.B) {
	bk := Backoff{Minimum: 100 * time.Millisecond, Maximum: 5, Jitter: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bk.duration(1)
	}
}

func ExampleBackoffDo(t *testing.T) {
	var conn net.Conn
	var err error

	// configure backoff, starting at 100ms up to 5 times (1.6s)
	b := New(100*time.Millisecond, 5)

	err = b.Do(func() error {
		conn, err = net.Dial("tcp", "exmaple.org:80")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalln("fatal! could not connect!", err)
	}
	defer conn.Close()

	// ... do something with conn ...
}

func TestBackoffDo(t *testing.T) {
	b := Backoff{Minimum: 10 * time.Millisecond, Maximum: 5}

	n := 0
	err := b.Do(func() error {
		n++
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, n, 1)
}

func TestBackoffDoError(t *testing.T) {
	b := Backoff{Minimum: 10 * time.Millisecond, Maximum: 5}

	n := 0
	testError := errors.New("[ test error ]")
	err := b.Do(func() error {
		n++
		return testError
	})
	assert.Equal(t, err, testError)
	assert.Equal(t, n, 5)
}

func TestBackoffDial(t *testing.T) {
	ts := tcptest.NewServer(func(conn net.Conn) {
		conn.Close()
	})
	defer ts.Close()

	b := Backoff{Minimum: 100 * time.Millisecond, Maximum: 5}

	conn, err := b.Dial("tcp", ts.Addr.String())
	if assert.NoError(t, err) {
		conn.Close()
	}
}
