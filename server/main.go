package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/aringq10/TCP/server/conn"
    "github.com/aringq10/TCP/server/match"
)

const PORT  = 6767
const MOUNT = "/ws"

func wsHandler(w http.ResponseWriter, r *http.Request) {
    c, err := conn.SetupConn(w, r)
    if err != nil {
        log.Println(err)
        return
    }

    log.Printf("WS: connection from %s, total %d", r.RemoteAddr, conn.Conns.Count())

    go c.ReadToChan()

    if p1, p2, ok := conn.Conns.TryDequeue2(); ok {
        go match.HandleMatch(p1, p2)
    }

    for {
        select {
        case _, ok := <-c.OutCh:
            // Read error or Conn closed before match start
            if !ok {
                conn.Conns.RemoveConn(c)
                log.Printf("WS: %s disconnected before match, %d", r.RemoteAddr, conn.Conns.Count())
                return
            }
        case <-c.DoneCh:
            // Match started
            log.Printf("WS: DoneCh signalled start of match for %s %d", r.RemoteAddr, conn.Conns.Count())
            return
        }
    }
}

func main() {
    address := fmt.Sprintf(":%d", PORT)

    log.Println("server listening on", PORT)

    http.HandleFunc(MOUNT, wsHandler)
    log.Fatal(http.ListenAndServe(address, nil))
}
