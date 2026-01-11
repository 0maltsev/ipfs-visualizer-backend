package kube

import (
	"fmt"

	"k8s.io/client-go/rest"
)

type ManualKubeConfigCreationError struct {
	Path  string
	Inner error
}

func (e ManualKubeConfigCreationError) Error() string {
	return fmt.Sprintf("kubernetes config with path %s not found or not formatted right", e.Path)
}

func (e ManualKubeConfigCreationError) Unwrap() error {
	return e.Inner
}

type AutoKubeConfigCreationError struct {
	Inner error
}

func (e AutoKubeConfigCreationError) Error() string {
	return fmt.Sprint("auto kubernetes config not found or not formatted right", e.Inner)
}

func (e AutoKubeConfigCreationError) Unwrap() error {
	return e.Inner
}

type ClientCreationError struct {
	Config *rest.Config
	Inner  error
}

func (e ClientCreationError) Error() string {
	return fmt.Sprintf("kubernetes client for config with host %s can't be created: %v", e.Config.Host, e.Inner)
}

func (e ClientCreationError) Unwrap() error {
	return e.Inner
}
