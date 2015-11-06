package backoff

import (
	"net"
	"testing"
	"time"

	"github.com/pusher/buddha/tcptest"
)

func TestBackoffDuration(t *testing.T) {
	b := New(100*time.Millisecond, 5)

	if d := b.Duration().String(); d != "100ms" {
		t.Fatal("expected 100ms, got", d)
	} else if d := b.Duration().String(); d != "200ms" {
		t.Fatal("expected 200ms, got", d)
	} else if d := b.Duration().String(); d != "400ms" {
		t.Fatal("expected 400ms, got", d)
	} else if d := b.Duration().String(); d != "800ms" {
		t.Fatal("expected 800ms, got", d)
	}

	if b.failures != 4 {
		t.Fatal("expected 4 failures, got", b.failures)
	}
}

func TestBackoffDurationJitter(t *testing.T) {
	b := New(100*time.Millisecond, 5)
	b.Jitter = true

	between(t, 100*time.Millisecond, 200*time.Millisecond, b.Duration())
	between(t, 100*time.Millisecond, 400*time.Millisecond, b.Duration())
}

func TestBackoffDurationInternal(t *testing.T) {
	b := New(500*time.Millisecond, 5)

	if d := b.duration(float64(500*time.Millisecond), 0).String(); d != "500ms" {
		t.Fatal("expected 500ms, got", d)
	} else if d := b.duration(float64(500*time.Millisecond), 1).String(); d != "1s" {
		t.Fatal("expected 1s, got", d)
	} else if d := b.duration(float64(500*time.Millisecond), 2).String(); d != "2s" {
		t.Fatal("expected 2s, got", d)
	} else if d := b.duration(float64(500*time.Millisecond), 3).String(); d != "4s" {
		t.Fatal("expected 4s, got", d)
	}
}

func TestBackoffDial(t *testing.T) {
	ts := tcptest.NewServer(func(conn net.Conn) {
		defer conn.Close()
	})
	defer ts.Close()

	b := New(100*time.Millisecond, 5)

	conn, err := b.Dial("tcp", ts.Addr.String())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer conn.Close()
}

func TestBackoffReset(t *testing.T) {
	b := New(100*time.Millisecond, 5)
	b.Duration()
	b.Duration()
	b.Duration()

	b.Reset()
	if b.failures != 0 {
		t.Fatal("expected 0 failures, got", b.failures)
	}
}

func BenchmarkBackoffDuration(b *testing.B) {
	bk := New(100*time.Millisecond, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bk.Duration()
	}
}

func between(t *testing.T, min, max, v time.Duration) {
	if v < min {
		t.Fatalf("expected > %s, got %s", min, v)
	}

	if v > max {
		t.Fatalf("expected < %s, got %s", max, v)
	}
}
