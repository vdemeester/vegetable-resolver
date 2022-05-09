package main

import (
	"context"
	"errors"
	"fmt"
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
		return errors.New(`require "kind" and "name" params`)
	}
	kind, hasKind := params["kind"]
	if !hasKind {
		return errors.New(`require "kind" param`)
	}
	kind = strings.TrimSpace(strings.ToLower(kind))
	if kind != "task" && kind != "pipeline" {
		return fmt.Errorf("unrecognized kind %q, only task and pipeline are supported", kind)
	}
	if _, has := params["name"]; !has {
		return errors.New(`require "name" param`)
	}
	return nil
}

// Resolve uses the given params to resolve the requested file or resource.
func (r *resolver) Resolve(ctx context.Context, params map[string]string) (framework.ResolvedResource, error) {
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
