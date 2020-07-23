package client

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	sensorpb "github.com/arshabbir/grpcstreamingclient/src/domain"

	"google.golang.org/grpc"
)

type client struct{}

type Client interface {
	Start()
	Recv(sensorpb.SensorServie_SensorPushClient, *sync.WaitGroup)
}

func NewClient() Client {
	return &client{}
}

func (c *client) Start() {

	servAddr := os.Getenv("SERVERADDR")

	if servAddr == "" {
		log.Println("Set the environment variable SERVERADDR")
		return
	}

	cc, err := grpc.Dial(servAddr, grpc.WithInsecure())

	if err != nil {
		log.Println("Error in grpc Dial")
		return
	}

	defer cc.Close()

	//senorpb.New

	client := sensorpb.NewSensorServieClient(cc)

	stream, err := client.SensorPush(context.Background())

	if err != nil {
		log.Println("Error in RPC")
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go c.Recv(stream, &wg)

	log.Println("Starting Sending .....")
	for {

		//For the request

		sensorID := fmt.Sprintf("%d", rand.Intn(10))
		reading := rand.Float32()
		requestID := fmt.Sprintf("%d", time.Now().Unix())

		log.Println(fmt.Sprintf("Sending..%s, %f, %s", sensorID, reading, requestID))

		if err := stream.Send(&sensorpb.SensorDataRequest{SensorId: sensorID, Reading: reading, RequestId: requestID}); err != nil {
			log.Println("Error in Sending to stream")
			break
		}
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
}

func (c *client) Recv(stream sensorpb.SensorServie_SensorPushClient, wg *sync.WaitGroup) {

	log.Println(" Receiving go routine spawned.....")

	for {

		mesg, err := stream.Recv()
		if err != nil {
			log.Println("Error receiving from client stream...")
			break
		}
		log.Println(fmt.Sprintf("Received : %s , %s , %d", mesg.GetRequestId(), mesg.GetSensorId(), mesg.GetStatus()))

	}

	wg.Done()

	return
}
