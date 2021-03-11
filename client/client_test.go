package client

import (
	"context"
	"log"
	"testing"

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const config = "/home/transwarp/.kube/config"
const testNs = "test"

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
	pfl := v1.PathFinderList{}
	err = pfclient.PathFinderV1(testNs).List(context.TODO(), &pfl, PathFinderListOption{})
	if err != nil {
		log.Println(err)
	}
	log.Println(pfl)

	pf := v1.PathFinder{}
	err = pfclient.PathFinderV1(testNs).Get(context.TODO(), "pathfinder-sample", &pf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pf)

	err = pfclient.PathFinderV1(testNs).GetByRegion(context.TODO(), "DEFAULT", &pf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pf)

	pf.Annotations["a"] = "b"
	err = pfclient.PathFinderV1(testNs).Update(context.TODO(), &pf, &client.UpdateOptions{})
	if err != nil {
		log.Fatal(err)
	}
}
