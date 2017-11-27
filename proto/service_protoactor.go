
package proto


import errors "errors"
import log "log"
import actor "github.com/AsynkronIT/protoactor-go/actor"
import remote "github.com/AsynkronIT/protoactor-go/remote"
import cluster "github.com/AsynkronIT/protoactor-go/cluster"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

	
var xPingServiceFactory func() PingService

func PingServiceFactory(factory func() PingService) {
	xPingServiceFactory = factory
}

func GetPingServiceGrain(id string) *PingServiceGrain {
	return &PingServiceGrain{ID: id}
}

type PingService interface {
	Init(id string)
		
	Ping(*PingRequest, cluster.GrainContext) (*PingResponse, error)
		
}
type PingServiceGrain struct {
	ID string
}

	
func (g *PingServiceGrain) Ping(r *PingRequest, options ...cluster.GrainCallOption) (*PingResponse, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*PingResponse, error) {
			pid, statusCode := cluster.Get(g.ID, "PingService")
			if statusCode != remote.ResponseStatusCodeOK {
				return nil, fmt.Errorf("Get PID failed with StatusCode: %v", statusCode)
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "Ping", MessageData: bytes}
			response, err := pid.RequestFuture(request, conf.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &PingResponse{}
				err = proto.Unmarshal(msg.MessageData, result)
				if err != nil {
					return nil, err
				}
				return result, nil
			case *cluster.GrainErrorResponse:
				return nil, errors.New(msg.Err)
			default:
				return nil, errors.New("Unknown response")
			}
		}
	
	var res *PingResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		} else {
			if conf.RetryAction != nil {
				conf.RetryAction(i)
			}
		}
	}
	return nil, err
}

func (g *PingServiceGrain) PingChan(r *PingRequest, options ...cluster.GrainCallOption) (<-chan *PingResponse, <-chan error) {
	c := make(chan *PingResponse)
	e := make(chan error)
	go func() {
		res, err := g.Ping(r, options...)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}
	

type PingServiceActor struct {
	inner PingService
}

func (a *PingServiceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xPingServiceFactory()
		id := ctx.Self().Id
		a.inner.Init(id[7:]) //skip "remote$"

	case actor.AutoReceiveMessage: //pass
	case actor.SystemMessage: //pass

	case *cluster.GrainRequest:
		switch msg.Method {
			
		case "Ping":
			req := &PingRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.Ping(req, ctx)
			if err == nil {
				bytes, err := proto.Marshal(r0)
				if err != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", err)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}
		
		}
	default:
		log.Printf("Unknown message %v", msg)
	}
}

	


//Why has this been removed?
//This should only be done on servers of the below Kinds
//Clients should not be forced to also be servers

//func init() {
//	
//	remote.Register("PingService", actor.FromProducer(func() actor.Actor {
//		return &PingServiceActor {}
//		})		)
//	
//}



// type pingService struct {
//	cluster.Grain
// }

// func (*pingService) Ping(r *PingRequest, cluster.GrainContext) (*PingResponse, error) {
// 	return &PingResponse{}, nil
// }



// func init() {
// 	//apply DI and setup logic

// 	PingServiceFactory(func() PingService { return &pingService{} })

// }





