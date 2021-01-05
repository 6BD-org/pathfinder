package common

import (
	"context"
	"reflect"

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"github.com/6BD-org/pathfinder/consts"
	"github.com/6BD-org/pathfinder/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListPathFinders(c client.Client, namespace string) (*v1.PathFinderList, error) {
	pl := v1.PathFinderList{}
	if err := c.List(context.TODO(), &pl, client.InNamespace(namespace)); err != nil {
		return nil, err
	}
	return &pl, nil
}

func GetPathFinderRegion(c client.Client, namespace string, region string) (*v1.PathFinder, error) {
	pl := v1.PathFinderList{}
	if err := c.List(context.TODO(), &pl, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	filted := utils.Filter(
		pl.Items,
		func(pf interface{}) bool { return pf.(v1.PathFinder).Spec.Region == region },
		reflect.TypeOf(v1.PathFinder{}),
	)

	if len(filted) > 1 {
		return nil, NewErr(consts.CODE_DUP_PF, consts.F_ERR_DUPLICATED_REGION, namespace, region)
	}
	if len(filted) < 1 {
		return nil, NewErr(consts.CODE_REGION_NOT_FOUND, consts.F_ERR_REGION_NOT_FOUND, namespace, region)
	}
	for _, pf := range pl.Items {
		if pf.Spec.Region == region {
			return &pf, nil
		}
	}
	return nil, NewErr(consts.CODE_REGION_NOT_FOUND, consts.F_ERR_REGION_NOT_FOUND, namespace, region)
}
