Backoff
=======

Simple implementation of an exponential backoff routine, primarily for use in network settings.

	go get github.com/jamescun/backoff

[![GoDoc](https://godoc.org/github.com/jamescun/backoff?status.svg)](https://godoc.org/github.com/jamescun/backoff)


Example
-------

```go
b := backoff.New(100 * time.Millisecond, 5)

conn, err := b.Dial("tcp", "flapping-network-service.local:8080")
if err != nil {
	// err is the last connection error
}
defer conn.Close()
```

License
-------

MIT
