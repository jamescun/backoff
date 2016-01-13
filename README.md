Backoff
=======

Simple implementation of an exponential backoff mechanism, including helpers for establishing network connections.

	go get github.com/jamescun/backoff

[![GoDoc](https://godoc.org/github.com/jamescun/backoff?status.svg)](https://godoc.org/github.com/jamescun/backoff)


Example
-------

```go
b := backoff.New(100 * time.Millisecond, 5)

err := b.Do(func() error {
	conn, err = net.Dial("tcp", "flapping-network-server.local:80")
	if err != nil {
		return err
	}

	return nil
})
```
