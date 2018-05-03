package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/satori/go.uuid"
)

//go:generate protoc --gofast_out=. topology.proto

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [number-of-nodes] [output-directory]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	numNodes, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
	}
	outputDir := os.Args[2]
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Directory \"%s\" does not exist!\n", outputDir)
		usage()
	}

	clusterId := uuid.NewV4().String()
	nodeIds := make([]string, numNodes)

	for i := 0; i < numNodes; i++ {
		nodeIds[i] = uuid.NewV4().String()
	}
	sort.Strings(nodeIds)

	top := &Topology{
		ClusterID: clusterId,
		NodeIDs:   nodeIds,
	}
	if buf, err := proto.Marshal(top); err != nil {
		log.Println(err)
		os.Exit(1)
	} else if err := ioutil.WriteFile(outputDir+"/topology", buf, 0666); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for i := 0; i < numNodes; i++ {
		if err := ioutil.WriteFile(outputDir+fmt.Sprintf("/node%d.id", i), []byte(nodeIds[i]), 0666); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully produced topology file and ID files for cluster with %d nodes.\n", numNodes)
	fmt.Printf("Files placed in \"%s\".\n", outputDir)
	fmt.Printf("Cluster ID: %s\n", clusterId)
	fmt.Printf("Node IDs: %v\n", nodeIds)
}
