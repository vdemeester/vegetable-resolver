package main

import (
	"context"

	"github.com/tektoncd/resolution/pkg/resolver/framework"
	"github.com/vdemeester/potato-resolver/pkg/resolver/potato"
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main("controller",
		framework.NewController(context.Background(), potato.New()),
	)
}
