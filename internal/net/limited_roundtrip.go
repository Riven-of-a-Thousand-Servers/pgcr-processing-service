package net

import (
	"io"
	"net/http"
)

type maxSizeTransport struct {
	MaxSize int64
	Base    http.RoundTripper
}

func (t *maxSizeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	resp.Body = &limitedReadCloser{
		r: io.LimitReader(resp.Body, t.MaxSize+1),
		c: resp.Body,
	}

	return resp, nil
}

type limitedReadCloser struct {
	r io.Reader
	c io.Closer
}

func (l *limitedReadCloser) Read(p []byte) (int, error) {
	return l.r.Read(p)
}

func (l *limitedReadCloser) Close() error {
	return l.c.Close()
}
