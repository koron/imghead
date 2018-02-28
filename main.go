package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	fetchSize int = 1024
	workerNum int = 4
	file      string
)

var l = log.New(os.Stderr, "IMGHEAD ", log.LstdFlags)

func inf2str(inf *ImageInfo) string {
	return fmt.Sprintf("statusCode:%d\tcontentLength:%d\twidth:%d\theight:%d\tformat:%s", inf.StatusCode, inf.ContentLength, inf.Width, inf.Height, inf.Format)
}

func exitCode(err error) int {
	switch err.(type) {
	case *FetchError:
		return 2
	case *DecodeError:
		return 3
	default:
		return 1
	}
}

func singleMode(u string) {
	inf, err := ImageHead(context.Background(), u, fetchSize)
	if err != nil {
		l.Printf("ERROR: %s", err)
		os.Exit(exitCode(err))
	}
	fmt.Println(inf2str(inf))
}

func multiMode(list []string) {
	ch, wg := startWorkers(context.Background(), workerNum)
	for _, u := range list {
		ch <- u
	}
	close(ch)
	wg.Wait()
}

func fileMode(r io.Reader) {
	rd := bufio.NewReader(r)
	ch, wg := startWorkers(context.Background(), workerNum)
	for {
		u, err := rd.ReadString('\n')
		if u != "" {
			ch <- strings.TrimSpace(u)
		}
		if err != nil {
			if err != io.EOF {
				l.Printf("failed to read input: %s", err)
			}
			break
		}
	}
	close(ch)
	wg.Wait()
}

func startWorkers(ctx context.Context, n int) (chan string, *sync.WaitGroup) {
	ch := make(chan string)
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			workerMain(i, ctx, ch)
			wg.Done()
		}(i)
	}
	return ch, wg
}

func workerMain(n int, ctx context.Context, ch chan string) {
	for {
		select {
		case u, ok := <-ch:
			if !ok {
				return
			}
			inf, err := ImageHead(ctx, u, fetchSize)
			if err != nil {
				l.Printf("failed for %s: %s", u, err)
				continue
			}
			fmt.Printf("%s\t%s\n", u, inf2str(inf))
		case _ = <-ctx.Done():
			return
		}
	}
}

func main() {
	flag.IntVar(&fetchSize, "size", 1024, "size to fetch")
	flag.IntVar(&workerNum, "worker", 4, "num of worker")
	flag.StringVar(&file, "file", "", "URL list file")
	flag.Parse()

	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			l.Fatalf("failed to read file: %s", err)
		}
		fileMode(f)
		f.Close()
		return
	}

	switch flag.NArg() {
	case 0:
		fileMode(os.Stdin)
	case 1:
		singleMode(flag.Args()[0])
	default:
		multiMode(flag.Args())
	}
}
