package main

import (
	"log"
	"net/http"

	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/vasiliy-t/patrain/proto"
)

func init() {
}

func main() {
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("mycluster", "172.19.0.10:9000", cp)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		g := proto.GetPingServiceGrain("actor")
		g.Ping(&proto.PingRequest{})
	})

	http.ListenAndServe("0.0.0.0:8080", nil)
}
