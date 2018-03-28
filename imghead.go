package main

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"

	_ "golang.org/x/image/bmp"
)

// ImageInfo is for ImageHead response.
type ImageInfo struct {
	StatusCode    int
	ContentLength int64

	Format string
	Width  int
	Height int
}

// ImageError is an error interface which occurred in imghead.
type ImageError interface {
	error
	StatusCode() int
}

// FetchError means HTTP fetch operation failure.
type FetchError struct {
	sc int
}

// StatusCode returns HTTP failed status code.
func (err *FetchError) StatusCode() int {
	return err.sc
}

// Error returns error message.
func (err *FetchError) Error() string {
	return fmt.Sprintf("fetch failed: status code=%d", err.sc)
}

// DecodeError means failure on decoding image.
type DecodeError struct {
	sc   int
	derr error
}

// StatusCode returns HTTP failed status code.
func (err *DecodeError) StatusCode() int {
	return err.sc
}

// Error returns error message.
func (err *DecodeError) Error() string {
	return "decode failed: " + err.derr.Error()
}

func imageHead(ctx context.Context, u string, n int) (*ImageInfo, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	if n > 0 {
		req.Header.Add("Range", "bytes=0-"+strconv.Itoa(n-1))
	}
	if ctx != nil {
		req.WithContext(ctx)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	sc := resp.StatusCode
	if sc != http.StatusOK && sc != http.StatusPartialContent {
		return nil, &FetchError{sc: sc}
	}
	c, f, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return nil, &DecodeError{sc: sc, derr: err}
	}
	return &ImageInfo{
		StatusCode:    sc,
		ContentLength: resp.ContentLength,
		Format:        f,
		Width:         c.Width,
		Height:        c.Height,
	}, nil
}

// ImageHead gets head of n bytes of URL, and decoded as image.
func ImageHead(ctx context.Context, u string, n int) (*ImageInfo, error) {
	// try partial content to decode.
	if n > 0 {
		inf, err := imageHead(ctx, u, n)
		if err == nil {
			return inf, nil
		}
		if _, ok := err.(ImageError); !ok {
			return nil, err
		}
		l.Printf("fallback to full content for %s: %s", u, err)
	}
	// decode full content.
	inf, err := imageHead(ctx, u, 0)
	if err != nil {
		return nil, err
	}
	return inf, nil
}
