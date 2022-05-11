package potato

#labels: [string]: string
#annotations: [string]: string
#param: {
	name: string
	description?: string
	default?: string | [...string]
	type: string | "array" | *"string"
}
#params: [...#param]
#result: {
	name: string
	description?: string
}
#workspace: {
	name: string
	description?: string
	optional?: bool
	mountPath?: string
}
#volumeMount: {
	name: string
	mountPath: string
}
#emptyDir: {}
#volume: {
	name: string
	emptyDir?: #emptyDir
}
#securityContext: {
	privileged?: bool
	runAsUser?: int
}
#env: {
	name: string
	valueFrom?: {
		secretKeyRef?: {
			name: string
			key: string
		}
	}
}

#metadata: {
	name:       string
	namespace?: string
	labels:     #labels
	annotations: #annotations
}

#step: {
	name?: string
	image: string
	command?: [...string]
	args?: [...string]
	env?: [...#env]
	workingDir?: string
	script?: string
	volumeMounts?: [...#volumeMount]
	securityContext?: #securityContext
}

#Task: {
	apiVersion: "tekton.dev/v1beta1"
	kind:       "Task"
	metadata: #metadata
	spec: {
		description?: string
		params?: [...#param]
		results?: [...#result]
		workspaces?: [...#workspace]
		steps?: [...#step]
		volumes?: [...#volume]
	}
}

#requiredAnnotations: #annotations & {
	"potato.tekton.dev/type":  string
}
#builderAnnotations: #requiredAnnotations & {
	"potato.tekton.dev/type": "builder"
}
#scanAnnotations: #requiredAnnotations & {
	"potato.tekton.dev/type": "scan"
}

#Builder: #Task & {
	metadata: #metadata & {
		annotations: #builderAnnotations
	}
	spec: {
		params: [{
			name: "IMAGE"
			type: string
		}, ...]
		workspaces: [{
			name: "source"
		}, ...]
		results: [{
			name:        "IMAGE_DIGEST"
		}, {
			name:        "IMAGE_URL"
		}, ...]
	}
}

#Scan: #Task & {
	metadata: #metadata & {
		annotations: #scanAnnotations
	}
	spec: {
		params: [{
			name: "IMAGE"
			type: string
		}, ...]
	}
}
