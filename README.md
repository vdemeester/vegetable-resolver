# vegetable resolvers

A set of experimental tektoncd/pipeline resolvers. The idea of this is
to experiment with what could be achieved with a Resolver.

## Potato 🥔 resolver


The main idea explored here is how to generate a Pipeline Spec
(embedded spec, or not) from a set of input, as well as, use "Task"
type to inject different types of Task into a Pipeline.

### Task type and injection

The idea behind *Task types* is more easily explained with an
example. Let's assume we have a "general" `Pipeline` that is building
an OCI image, pushing it on some registry, scan the content of the
image and run tests against it.  There is *a lot* of options to build
an image, and most likely, we don't necessarily care about which tool
is used.

This means, a Pipeline could *specify* that "an image needs to be
built" in that task, and from there we could "inject" a Task that does
this.

To keep the example, a `Task` that builds an image should "fulfil" the
following contract:

```yaml
# This is a partial Task definition with the only
# required field.
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  annotations:
    potato.tekton.dev/type: builder
spec:
  params:
  - name: IMAGE
    description: Name (reference) of the image to build
  results:
  - name: IMAGE_DIGEST
    description: Digest of the image built
  - name: IMAGE_URL
    description: Digest of the image built
```

Let's "name" this contract `builder`.

A pipeline could "ask" for this `builder` contract, like below.

```yaml
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: image-pipeline
  annotations:
    # Default builder to use is kaniko
    default.potato.tekton.dev/builder: kaniko
spec:
  workspaces:
  - name: shared-workspace
  tasks:
  - name: checkout
    taskRef:
      name: git-clone
    workspaces:
    - name: output
      workspace: shared-workspace
  - name: build-and-push
    taskRef:
      # Any Task that fulfil the "builder" contract
      name: {{ builder }}
    params:
    - name: IMAGE
      value: quay.io/vdemeest/foo:bar
    workspaces:
    - name: source
      workspace: shared-workspace
  - name: scan
    params:
    - name: IMAGE
      value: $(tasks.build-and-push.results.IMAGE_URL)
```

This `potato` resolver would then allow the user to inject its own
`builder` Task or rely on a configured default.

```yaml
# Relying on the default build
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  generateName: build-me-an-image
spec:
  pipelineRef:
    resolver: potato
    resource:
    - name: pipeline
      value: image-pipeline
---
# Or specifying our own
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  generateName: build-me-an-image
spec:
  pipelineRef:
    resolver: potato
    resource:
    - name: pipeline
      value: image-pipeline
    - name: builder
      value: buildah
```
