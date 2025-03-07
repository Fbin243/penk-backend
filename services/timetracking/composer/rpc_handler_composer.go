package composer

import "tenkhours/services/timetracking/transport/rpc"

func ComposeRPCHandler() *rpc.TimetrackingHandler {
	composer = GetComposer()
	return rpc.NewtimetrackingHandler(composer.timetrackingBiz)
}
