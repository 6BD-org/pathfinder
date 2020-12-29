package client

import (
	"context"
	"log"
	"testing"

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const config = "/root/.kube/config"
const testNs = "makaveli"

func Test(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		log.Println(err)
		return
	}
	pfclient, err := New(config)
	if err != nil {
		log.Println(err)
		return
	}
	pf := v1.PathFinderList{}
	err = pfclient.PathFinderV1(testNs).List(context.TODO(), &pf, PathFinderListOption{Region: "DEFAULT"})
	if err != nil {
		log.Println(err)
	}
	log.Println(pf)

}
