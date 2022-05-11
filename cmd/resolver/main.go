package main

import (
	"github.com/tektoncd/resolution/pkg/resolver/framework"
	"github.com/vdemeester/vegetable-resolver/pkg/resolver/carrot"
	"github.com/vdemeester/vegetable-resolver/pkg/resolver/potato"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
)

func main() {
	ctx := signals.NewContext()
	sharedmain.Main("controller",
		framework.NewController(ctx, potato.New()),
		framework.NewController(ctx, carrot.New()),
	)
}
