package potato

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelinesclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	"github.com/tektoncd/pipeline/pkg/resolution/common"
	"github.com/tektoncd/pipeline/pkg/resolution/resolver/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "knative.dev/pkg/logging"
)

const jsonContentType = "application/json"

func New() framework.Resolver {
	return &resolver{}
}

type resolver struct {
	// The clientset used to look up tasks and pipelines from the
	// clusterresolver's private namespace.
	Pipelineclientset pipelinesclientset.Interface
}

// Initialize creates an instance of the pipelines clientset so that
// tasks and pipelines can be looked up.
func (r *resolver) Initialize(ctx context.Context) error {
	r.Pipelineclientset = pipelineclient.Get(ctx)
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
func (r *resolver) ValidateParams(ctx context.Context, params []v1.Param) error {
	if len(params) == 0 {
		return errors.New(`require at least "pipeline" param`)
	}
	hasPipeline := false
	for _, p := range params {
		if p.Name == "pipeline" {
			hasPipeline = true
		}
	}
	if !hasPipeline {
		return errors.New(`require "pipeline" param`)
	}
	return nil
}

// Resolve uses the given params to resolve the requested file or resource.
func (r *resolver) Resolve(ctx context.Context, params []v1.Param) (framework.ResolvedResource, error) {
	// logger := logging.FromContext(ctx)
	pipeline := ""
	builder := ""
	for _, p := range params {
		if p.Name == "pipeline" {
			pipeline = p.Value.StringVal
		}
		if p.Name == "builder" {
			builder = p.Value.StringVal
		}
	}
	namespace := common.RequestNamespace(ctx)
	resolved, err := r.Pipelineclientset.TektonV1beta1().Pipelines(namespace).Get(ctx, pipeline, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error fetching pipeline %q", pipeline)
	}
	// FIXME: support dynamic contracts, …
	if builder != "" {
		// Read the annotation
		builder = resolved.Annotations["default.potato.tekton.dev/builder"]
	}
	// Manually add type meta because the kube api doesn't
	// necessarily include them in its response.
	out := resolved.DeepCopy()
	out.TypeMeta.Kind = "Pipeline"
	out.TypeMeta.APIVersion = "tekton.dev/v1beta1"

	for i, t := range resolved.Spec.Tasks {
		// FIXME: support dynamic contracts
		if t.TaskRef != nil && t.TaskRef.Name == "potato.type.builder" {
			nt := t.DeepCopy()
			nt.TaskRef.Name = builder
			out.Spec.Tasks[i] = *nt
		}
	}

	data, err := json.Marshal(*out)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal resolved resource to json: %w", err)
	}
	return &resolvedResource{
		data: data,
		source: &v1.RefSource{
			URI: "potato://foo-is-bar",
		},
	}, nil
}

// resolvedResource wraps the data we want to return to Pipelines
type resolvedResource struct {
	data   []byte
	source *v1.RefSource
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

func (r *resolvedResource) RefSource() *v1.RefSource {
	return r.source
}
