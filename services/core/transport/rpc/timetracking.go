package rpc

// func (hdl *CoreHandler) CreateTimeTracking(ctx context.Context, req *core.CreateTimeTrackingReq) (*core.TimeTracking, error) {
// 	// Map RPC input to entity input
// 	timeTrackingInput, err := Map[core.CreateTimeTrackingReq, entity.TimeTrackingInput](req, UnixTimeConverter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Call business logic to create time tracking
// 	timeTracking, err := hdl.timetrackingBiz.CreateTimeTracking(ctx, timeTrackingInput)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Map entity to RPC response
// 	return Map[entity.TimeTracking, core.TimeTracking](timeTracking, nil)
// }

// func (hdl *CoreHandler) UpdateTimeTracking(ctx context.Context, req *common.EmptyReq) (*core.TimeTracking, error) {
// 	// Call business logic to update time tracking
// 	timeTracking, err := hdl.timetrackingBiz.UpdateTimeTracking(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Map entity to RPC response
// 	return Map[entity.TimeTracking, core.TimeTracking](timeTracking, UnixTimeConverter)
// }

// // TODO: @Fbin243 implements the following methods later
// func (hdl *CoreHandler) GetCurrentTimeTracking(ctx context.Context, req *common.EmptyReq) (*core.TimeTracking, error) {
// 	return nil, errors.New("not implemented")
// }

// func (hdl *CoreHandler) GetTotalCurrentTimeTracking(ctx context.Context, req *core.TotalTimeTrackingReq) (*core.TotalTimeTrackingResp, error) {
// 	return nil, errors.New("not implemented")
// }
