package main

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"

	_ "golang.org/x/image/bmp"
)

// ImageInfo is for ImageHead response.
type ImageInfo struct {
	StatusCode int

	Format string
	Width  int
	Height int
}

type ImageError interface {
	error
	StatusCode() int
}

type imageError struct {
	sc   int
	derr error
}

func (err *imageError) StatusCode() int {
	return err.sc
}

func (err *imageError) Error() string {
	if err.derr != nil {
		return "decode error: " + err.derr.Error()
	}
	return fmt.Sprintf("no contents: status code is %d", err.sc)
}

func imageHead(ctx context.Context, u string, n int) (*ImageInfo, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	if n > 0 {
		req.Header.Add("Range", "bytes=0-%d"+strconv.Itoa(n))
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
		return nil, &imageError{sc: sc}
	}
	c, f, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return nil, &imageError{sc: sc, derr: err}
	}
	return &ImageInfo{
		StatusCode: sc,
		Format:     f,
		Width:      c.Width,
		Height:     c.Height,
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
		log.Printf("fallback to full content for %s: %s", u, err)
	}
	// decode full content.
	inf, err := imageHead(ctx, u, 0)
	if err != nil {
		return nil, err
	}
	return inf, nil
}
