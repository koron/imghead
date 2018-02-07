package main

import (
	"context"
	"flag"
	"fmt"
	"log"
)

var (
	fetchSize int = 1024
)

func inf2str(inf *ImageInfo) string {
	return fmt.Sprintf("statusCode:%d\twidth:%d\theight:%d\tformat:%s", inf.StatusCode, inf.Width, inf.Height, inf.Format)
}

func singleMode(u string) {
	inf, err := ImageHead(context.Background(), u, 1024)
	if err != nil {
		log.Fatal("ERROR: %s", err)
	}
	fmt.Println(inf2str(inf))
}

func multiMode(list []string) {
	// TODO:
	log.Fatal("multi mode is not implemented")
}

func fileMode() {
	// TODO:
	log.Fatal("file mode is not implemented")
}

func main() {
	flag.IntVar(&fetchSize, "size", 1024, "size to fetch")
	flag.Parse()
	switch flag.NArg() {
	case 0:
		fileMode()
	case 1:
		singleMode(flag.Args()[0])
	default:
		multiMode(flag.Args())
	}
}
