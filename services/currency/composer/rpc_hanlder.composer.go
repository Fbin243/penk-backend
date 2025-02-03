package composer

import "tenkhours/services/currency/transport/rpc"

func ComposeRPCHandler() *rpc.CurrencyHandler {
	composer = GetComposer()
	return rpc.NewCurrencyHandler(composer.CurrencyBiz)
}
