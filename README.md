# imghead tool

imghead fetch range of image header and decode its dimension.
if the server doesn't support range requests, it fetch whole of the content.

## Install or update

Install or update:

```console
$ go get -u github.com/koron/imghead
```

## Getting started

Fetch single image header

```console
$ imghead.exe http://httpbin.org/image/png
statusCode:200  contentLength:8090      width:100       height:100      format:png
```

Fetch header of multiple images by arguments

```console
$ imghead.exe http://httpbin.org/image/png http://httpbin.org/image/jpeg
http://httpbin.org/image/png    statusCode:200  contentLength:8090      width:100       height:100      format:png
http://httpbin.org/image/jpeg   statusCode:200  contentLength:35588     width:239       height:178      format:jpeg
```

Fetch header of multiple images by file

```console
$ cat list.txt
http://httpbin.org/image/png
http://httpbin.org/image/jpeg

$ imghead -file list.txt
http://httpbin.org/image/png    statusCode:200  contentLength:8090      width:100       height:100      format:png
http://httpbin.org/image/jpeg   statusCode:200  contentLength:35588     width:239       height:178      format:jpeg
```

### Options

*   `-size` - size of range query (default: 1024)
*   `-worker` - number of parallel download
*   `-file` - file for URL list to fetch

## Input URLs priority

1. `-file` option
2. arguments
3. STDIN
