package composer

import "tenkhours/services/analytic/transport/rpc"

func ComposeRPCHandler() *rpc.AnalyticHandler {
	composer = GetComposer()
	return rpc.NewAnalyticHandler(composer.AnalyticBiz)
}
