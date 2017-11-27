package main

import (
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/vasiliy-t/patrain/proto"
)

func init() {
	remote.Register("PingService", actor.FromProducer(func() actor.Actor {
		return &proto.PingServiceActor{}
	}))
}

func main() {
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("mycluster", "172.19.0.11:9001", cp)
	<-time.After(1 * time.Hour)
}
