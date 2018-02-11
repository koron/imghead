# imghead tool

imghead fetch head of image and decode its dimension.

## Input URLs priority

1. `-file` option
2. arguments
3. STDIN

## Examples

Fetch single image header

```console
$ imghead https://example.org/image001.png
statusCode:200	width:1200	height:675	format:png
```

Fetch header of multiple images by arguments

```console
$ imghead https://example.org/image001.png https://example.org/image002.png
https://example.org/image001.png	statusCode:200	width:1200	height:675	format:png
https://example.org/image002.png	statusCode:200	width:800	height:600	format:png
```

Fetch header of multiple images by file

```console
$ cat list.txt
https://example.org/image001.png
https://example.org/image002.png
https://example.org/image003.png
https://example.org/image004.png
https://example.org/image005.png

$ imghead -file list.txt
https://example.org/image001.png	statusCode:200	width:1200	height:675	format:png
https://example.org/image002.png	statusCode:200	width:800	height:600	format:png
https://example.org/image003.png	statusCode:200	width:800	height:600	format:png
https://example.org/image004.png	statusCode:200	width:800	height:600	format:png
https://example.org/image005.png	statusCode:200	width:800	height:600	format:png
```

### Options

*   `-size` - size of range query (default: 1024)
*   `-worker` - number of parallel download
*   `-file` - file for URL list to fetch
