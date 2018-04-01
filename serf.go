package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "time"

	// "github.com/hashicorp/memberlist"
	// "github.com/hashicorp/raft"
	// "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
)

type SerfAgent struct {
	FinishChan  chan int
	EventsChan  chan serf.Event
	config      *serf.Config
	agent       *serf.Serf
	clusterName string
	wg          *sync.WaitGroup
}

func NewSerfAgent(nodeName string, clusterName string) *SerfAgent {
	eventsChan := make(chan serf.Event, 16)

	config := serf.DefaultConfig()
	config.NodeName = nodeName
	config.EventCh = eventsChan
	config.LogOutput = os.Stdout

	agent, err := serf.Create(config)
	if err != nil {
		log.Fatal(err)
	}

	return &SerfAgent{
		FinishChan:  make(chan int, 1),
		EventsChan:  eventsChan,
		config:      config,
		agent:       agent,
		clusterName: clusterName,
		wg:          new(sync.WaitGroup),
	}
}

func (sa *SerfAgent) Run() {
	// Join an existing cluster by specifying at least one known member.
	sa.agent.Join([]string{sa.clusterName}, false)
	sa.wg.Add(1)

LOOP:
	for {
		select {
		case <-sa.FinishChan:
			if err := sa.agent.Leave(); err != nil {
				fmt.Printf("Error while leave Serf Cluster. %v\n", err)
			}
			break LOOP
		case event := <-sa.EventsChan:
			fmt.Printf("An event arrived: %v", event)
		}
	}

	fmt.Printf("Serf Agent will stop. Leaving from cluster\n")
	sa.wg.Done()
}

func (sa *SerfAgent) Join() {
	sa.wg.Wait()
}

func main() {
	nodeIp := os.Getenv("NODE_ADDR")
	clusterIp := os.Getenv("CLUSTER_ADDR")
	nodePort := os.Getenv("NODE_PORT")
	clusterPort := os.Getenv("CLUSTER_PORT")

	nodeName := fmt.Sprintf("%s:%s", nodeIp, nodePort)
	clusterName := fmt.Sprintf("%s:%s", clusterIp, clusterPort)

	serfAgent := NewSerfAgent(nodeName, clusterName)
	go serfAgent.Run()

	handleSigintSignal()

	// Wait the SerfAgent thread to finish before continue
	serfAgent.FinishChan <- 1
	serfAgent.Join()
}

func handleSigintSignal() {
	c := make(chan os.Signal, syscall.SIGINT)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Printf("[MAIN] SIGNIT signal received. Stopping program.\n")
}
