package auth

import (
	"reflect"

	"tenkhours/pkg/utils"

	"github.com/casbin/casbin/v2"
)

var enforcer *casbin.Enforcer

func init() {
	modelPath := utils.GetRoot() + "/pkg/auth/abac_model.conf"
	policyPath := utils.GetRoot() + "/pkg/auth/abac_policy.csv"

	var err error
	enforcer, err = casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		panic(err)
	}
}

func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

func CheckPermission(sub, obj, act any) (bool, error) {
	subTemp := sub
	objTemp := obj
	if reflect.TypeOf(sub).Kind() == reflect.Ptr {
		subTemp = reflect.ValueOf(sub).Elem().Interface()
	}

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		objTemp = reflect.ValueOf(obj).Elem().Interface()
	}

	return enforcer.Enforce(reflect.TypeOf(subTemp).Name(), reflect.TypeOf(objTemp).Name(), subTemp, objTemp, act)
}
