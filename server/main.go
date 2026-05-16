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
    _, err := conn.SetupConn(w, r)
    if err != nil {
        log.Println(err)
        return
    }

    log.Printf("ws connection from %s, total %d", r.RemoteAddr, conn.Conns.Count())

    if conn.Conns.Count() >= 2 {
        go match.HandleMatch(conn.Conns.Get(0), conn.Conns.Get(1))
    }
}

func main() {
    address := fmt.Sprintf(":%d", PORT)

    log.Println("server listening on", PORT)

    http.HandleFunc(MOUNT, wsHandler)
    log.Fatal(http.ListenAndServe(address, nil))
}
