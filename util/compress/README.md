# UnpackIt

This Go library allows you to easily unpack the following files using magic numbers:

* tar.gz
* tar.bzip2
* tar.xz
* zip
* tar

## Usage

Unpack a file:

```go
    file, _ := os.Open(test.filepath)
    destPath, err := gmccompress.Unpack(file, tempDir)
```

Unpack a stream (such as a http.Response):

```go
    res, err := http.Get(url)
    destPath, err := gmccompress.Unpack(res.Body, tempDir)
```

