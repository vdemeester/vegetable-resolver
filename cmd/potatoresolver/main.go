package main

import (
	"context"
	"errors"
	"fmt"
	// "strconv"
	"strings"

	"github.com/tektoncd/resolution/pkg/common"
	"github.com/tektoncd/resolution/pkg/resolver/framework"
	"knative.dev/pkg/injection/sharedmain"
)

const clusterResolverPrivateNamespace = "tekton-cluster-scoped-resources"
const jsonContentType = "application/json"

func main() {
	sharedmain.Main("controller",
		framework.NewController(context.Background(), &resolver{}),
	)
}

type resolver struct {
	// The clientset used to look up tasks and pipelines from the
	// clusterresolver's private namespace.
	// Pipelineclientset pipelinesclientset.Interface
}

// Initialize creates an instance of the pipelines clientset so that
// tasks and pipelines can be looked up.
func (r *resolver) Initialize(ctx context.Context) error {
	// r.Pipelineclientset = pipelineclient.Get(ctx)
	return nil
}

// GetName returns a string name to refer to this resolver by.
func (r *resolver) GetName(context.Context) string {
	return "Potato"
}

// GetSelector returns a map of labels to match requests to this resolver.
func (r *resolver) GetSelector(context.Context) map[string]string {
	return map[string]string{
		common.LabelKeyResolverType: "potato",
	}
}

// ValidateParams ensures parameters from a request are as expected.
// Only "kind" and "name" are needed.
func (r *resolver) ValidateParams(ctx context.Context, params map[string]string) error {
	if len(params) == 0 {
		return errors.New(`require at least "type" param`)
	}
	t, hasType := params["type"]
	if !hasType {
		return errors.New(`require "type" param`)
	}
	t = strings.TrimSpace(strings.ToLower(t))
	// FIXME change this
	if t != "dockerfile" {
		return fmt.Errorf("unrecognized type %q, only dockerfile is supported", t)
	}
	return nil
}

// Resolve uses the given params to resolve the requested file or resource.
func (r *resolver) Resolve(ctx context.Context, params map[string]string) (framework.ResolvedResource, error) {
	t := params["type"]
	switch t {
	case "dockerfile":
		// imagebuilder := params["imagebuilder"]
		// scan, err := strconv.ParseBool(params["scan"])
		// if err != nil {
		// return nil, err
		// }
		return &resolvedResource{
			data: []byte(`---
apiVersion: tekton.dev/v1beta1
kind: Pipeline
spec:
  workspaces:
  - name: shared-workspace
  - name: sslcertdir
    optional: true
  tasks:
  - name: fetch-repository
    taskRef:
      name: git-clone
    workspaces:
    - name: output
      workspace: shared-workspace
    params:
    - name: url
      value: https://github.com/kelseyhightower/nocode
    - name: subdirectory
      value: ""
    - name: deleteExisting
      value: "true"
  - name: buildah
    taskRef:
      name: buildah
    runAfter:
    - fetch-repository
    workspaces:
    - name: source
      workspace: shared-workspace
    - name: sslcertdir
      workspace: sslcertdir
    params:
    - name: IMAGE
      value: registry:5000/nocode
`),
		}, nil
	}
	return nil, fmt.Errorf("not implemented yet")
}

// resolvedResource wraps the data we want to return to Pipelines
type resolvedResource struct {
	data []byte
}

// Data returns the bytes of the task or pipeline resolved from the
// private namespace.
func (r *resolvedResource) Data() []byte {
	return r.data
}

// Annotations returns a content-type of json since the data is
// serialized as json.
func (r *resolvedResource) Annotations() map[string]string {
	return map[string]string{
		common.AnnotationKeyContentType: jsonContentType,
	}
}
