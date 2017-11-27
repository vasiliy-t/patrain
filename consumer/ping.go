package main

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/vasiliy-t/patrain/proto"
)

type pingService struct {
	cluster.Grain
	counter int64
}

func (s *pingService) Ping(r *proto.PingRequest, ctx cluster.GrainContext) (*proto.PingResponse, error) {
	s.counter++
	log.Printf("COUNTER %d", s.counter)
	return &proto.PingResponse{}, nil
}

func init() {
	proto.PingServiceFactory(func() proto.PingService { return &pingService{} })
}
