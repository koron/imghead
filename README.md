# imghead tool

imghead fetch range of image header and decode its dimension.
if the server doesn't support range requests, it fetch whole of the content.

## Install or update

Install or update:

```console
$ go get -u github.com/koron/imghead
```

## Getting started

Fetch single image header.  If it failed, reason of failure is known by its
exit code.  See [below](#exit-code) for details.

```console
$ imghead https://pbs.twimg.com/media/DU3Jo_LUMAA0Zpw.jpg
statusCode:206  contentLength:1024      width:900       height:1200     format:jpeg
```

Fetch single image by arguments

```console
$ imghead http://httpbin.org/image/png https://pbs.twimg.com/media/DU3Jo_LUMAA0Zpw.jpg
http://httpbin.org/image/png    statusCode:200  contentLength:8090      width:100       height:100      format:png
https://pbs.twimg.com/media/DU3Jo_LUMAA0Zpw.jpg statusCode:206  contentLength:1024      width:900       height:1200     format:jpeg
```

Fetch header of multiple images by file

```console
$ cat testdata/list.txt
https://pbs.twimg.com/media/DVuoPdOV4AEA_rV.jpg
https://pbs.twimg.com/media/DVlbdZ3VQAADJMl.jpg
https://pbs.twimg.com/media/DVkLilJVQAAH1h9.jpg
https://pbs.twimg.com/media/DVab9nmX0AEro8f.jpg
https://pbs.twimg.com/media/DU3Jo_LUMAA0Zpw.jpg

$ imghead -file testdata/list.txt
https://pbs.twimg.com/media/DVab9nmX0AEro8f.jpg statusCode:206  contentLength:1024      width:1200      height:900      format:jpeg
https://pbs.twimg.com/media/DVkLilJVQAAH1h9.jpg statusCode:206  contentLength:1024      width:1200      height:900      format:jpeg
https://pbs.twimg.com/media/DVuoPdOV4AEA_rV.jpg statusCode:206  contentLength:1024      width:900       height:1200     format:jpeg
https://pbs.twimg.com/media/DVlbdZ3VQAADJMl.jpg statusCode:206  contentLength:1024      width:1200      height:900      format:jpeg
https://pbs.twimg.com/media/DU3Jo_LUMAA0Zpw.jpg statusCode:206  contentLength:1024      width:900       height:1200     format:jpeg
```

### Options

*   `-size` - size of range query (default: 1024)
*   `-worker` - number of parallel download
*   `-file` - file for URL list to fetch

## Input URLs priority

1.  `-file` option
2.  arguments
3.  STDIN

## Supported image formats

*   GIF
*   PNG
*   JPEG
*   BMP

## Exit code

imghead returns below exit codes when it failed with single argument.

*   failed to fetch HTTP/S: 2
*   failed to decode image: 3
*   other failure: 1
