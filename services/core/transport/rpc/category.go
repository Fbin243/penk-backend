package rpc

import (
	"context"

	"tenkhours/proto/pb/common"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertCategory(ctx context.Context, req *core.CategoryInput) (*core.Category, error) {
	categoryInput, err := Map[core.CategoryInput, entity.CategoryInput](req, UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	category, err := hdl.categoryBiz.Upsert(ctx, categoryInput)
	if err != nil {
		return nil, err
	}

	return Map[entity.Category, core.Category](category, UnixTimeConverter)
}

func (hdl *CoreHandler) DeleteCategory(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	category, err := hdl.categoryBiz.Delete(ctx, req.Id)

	return &common.IdResp{
		Id: category.ID,
	}, err
}
