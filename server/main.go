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

    go c.ReadToChan()

    if p1, p2, ok := conn.Conns.TryDequeue2(); ok {
        go match.HandleMatch(p1, p2)
    }

    for {
        select {
        case _, ok := <-c.OutCh:
            // Read error or Conn closed before match start
            // Otherwise, discard messages
            if !ok {
                conn.Conns.RemoveConn(c)
                return
            }
        case <-c.MatchStartCh:
            // Match started: stop discarding messages
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
