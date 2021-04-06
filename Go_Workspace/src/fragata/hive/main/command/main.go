
package main

import (
      "os"
      "io"
      "flag"
      "fmt"
      "fragata/hive/client"
)

func main() {
    scheme := flag.String("scheme", "http", "Hive server URL scheme")
    host := flag.String("host", "", "Hive server host name")
    port := flag.String("port", "", "Hive server port number")

    cookieDir := flag.String("cookie-dir", "", "Cookie directory")

    from := flag.String("from", "", "Query parameter 'from'")
    size := flag.String("size", "", "Query parameter 'size'")
    sortBy := flag.String("sort-by", "", "Query parameter 'sortBy'")
    sortDir := flag.String("sort-dir", "", "Query parameter 'sortDir'")
    task := flag.String("task", "", "Query parameter 'task'")
    state := flag.String("state", "", "Query parameter 'state'")
    verified := flag.String("verified", "", "Query parameter 'verified'")

    flag.Parse()

    aScheme := *scheme
    aHost := *host
    aPort := *port
    aCookieDir := *cookieDir
    if aHost == "" {
        aHost = os.Getenv("HIVE_HOST")
    }
    if aPort == "" {
        aPort = os.Getenv("HIVE_PORT")
    }
    if aCookieDir == "" {
        aCookieDir = os.Getenv("HIVE_COOKIE_DIR")
    }
    if aHost == "" {
        fmt.Fprintf(os.Stderr, "Error: Missing Hive host name\n")
        os.Exit(1)
    }

    w := os.Stdout    // Unix version
//    w = makeDosWriter(w)    // Window version; use raw os.Stdout on Unix
    clnt := client.NewClient(aScheme, aHost, aPort)
    cmd := client.NewCommand(clnt, w, aCookieDir, client.RoleAdmin|client.RoleUser)
    params := &client.Params{
        From: *from,
        Size: *size,
        SortBy: *sortBy,
        SortDir: *sortDir,
        Task: *task,
        State: *state,
        Verified: *verified,
    }
    err := cmd.Do(params, flag.Args())
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
        os.Exit(1)
    }
}

type dosWriter struct {
    w io.Writer
    buf []byte
}

func makeDosWriter(w io.Writer) *dosWriter {
    return &dosWriter{w, make([]byte, 4096)}
}

func(self *dosWriter) Write(p []byte) (int, error) {
    self.buf = self.buf[:0]
    lenp := len(p)
    for i := 0; i < lenp; i++ {
        c := p[i]
        if c == '\n' {
            self.buf = append(self.buf, '\r')
        }
        self.buf = append(self.buf, c)
    }
    return self.w.Write(self.buf)
}

