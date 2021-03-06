# Go Utils

Utils for Go dev.

## Usage

### 生成自签名 CA 证书

`生成证书`

```go
func main(){
    ip := []byte("127.0.0.1")
    alternateDNS := []string{"localhost"}
    GenerateSelfSignedCertKey("/root", 2048, "192.168.0.1", []net.IP{net.ParseIP(string(ip))}, alternateDNS)
}
```

`查看证书`

```bash
openssl x509 -in ca.crt -noout -text
```

### WriteCounter 使用

See <https://golangcode.com/download-a-file-with-progress/>

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "github.com/designinlife/goutils"
)

func main() {
    fileUrl := "https://upload.wikimedia.org/wikipedia/commons/d/d6/Wp-w4-big.jpg"
    err := DownloadFile("avatar.jpg", fileUrl)
    if err != nil {
        panic(err)
    }
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string) error {
    // Create the file, but give it a tmp file extension, this means we won't overwrite a
    // file until it's downloaded, but we'll remove the tmp extension once downloaded.
    out, err := os.Create(filepath + ".tmp")
    if err != nil {
        return err
    }

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        out.Close()
        return err
    }
    defer resp.Body.Close()

    // Create our progress reporter and pass it to be used alongside our writer
    counter := &WriteCounter{}
    if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
        out.Close()
        return err
    }

    // The progress use the same line so print a new line once it's finished downloading
    fmt.Print("\n")

    // Close the file without defer so it can happen before Rename()
    out.Close()

    if err = os.Rename(filepath+".tmp", filepath); err != nil {
        return err
    }
    return nil
}
```
