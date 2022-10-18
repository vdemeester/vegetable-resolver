package main

import (
	"context"

	"github.com/tektoncd/pipeline/pkg/resolution/resolver/framework"
	"github.com/vdemeester/vegetable-resolver/pkg/resolver/carrot"
	"github.com/vdemeester/vegetable-resolver/pkg/resolver/potato"
	"knative.dev/pkg/injection/sharedmain"
	// "knative.dev/pkg/signals"
)

func main() {
	// ctx := signals.NewContext()
	ctx := context.Background()
	sharedmain.Main("controller",
		framework.NewController(ctx, potato.New()),
		framework.NewController(ctx, carrot.New()),
	)
}
