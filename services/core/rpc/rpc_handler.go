package rpc

import (
	"context"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/pb"
	"tenkhours/services/core/business"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RPCHandler struct {
	pb.UnimplementedCoreServer
	profilesBusiness   *business.ProfilesBusiness
	charactersBusiness *business.CharactersBusiness
}

func NewRPCHandler(profilesBusiness *business.ProfilesBusiness, charactersBusiness *business.CharactersBusiness) *RPCHandler {
	return &RPCHandler{
		profilesBusiness:   profilesBusiness,
		charactersBusiness: charactersBusiness,
	}
}

func (s *RPCHandler) UpdateTimeInCharacter(ctx context.Context, req *pb.UpdateTimeReq) (*pb.UpdateTimeResp, error) {
	resp := &pb.UpdateTimeResp{Success: false}
	characterID, err := primitive.ObjectIDFromHex(req.CharacterID)
	if err != nil {
		return resp, err
	}

	metricID := primitive.NilObjectID
	if req.MetricID != "" {
		metricID, err = primitive.ObjectIDFromHex(req.MetricID)
		if err != nil {
			return resp, err
		}
	}

	err = s.charactersBusiness.UpdateTimeInCharacter(ctx, characterID, metricID, req.Time)
	if err != nil {
		return resp, err
	}

	resp.Success = true

	return resp, nil
}

func (s *RPCHandler) IntrospectProfile(ctx context.Context, req *pb.IntrospectReq) (*pb.IntrospectResp, error) {
	resp := &pb.IntrospectResp{Success: false}

	firebaseProfile := auth.FirebaseProfile{
		UID:     req.FirebaseUID,
		Email:   req.Email,
		Name:    req.Name,
		Picture: req.Picture,
	}

	profile, err := s.profilesBusiness.IntrospectProfile(ctx, firebaseProfile)
	if err != nil {
		return resp, err
	}

	resp.Success = true
	resp.ProfileID = profile.ID.Hex()

	return resp, nil
}

func (s *RPCHandler) CheckPermission(ctx context.Context, req *pb.CheckPermissionReq) (*pb.CheckPermissionResp, error) {
	resp := &pb.CheckPermissionResp{Authorized: false}

	profileID, err := primitive.ObjectIDFromHex(req.ProfileID)
	if err != nil {
		return resp, err
	}

	characterID, err := primitive.ObjectIDFromHex(req.CharacterID)
	if err != nil {
		return resp, err
	}

	metricID, err := primitive.ObjectIDFromHex(req.MetricID)
	if err != nil {
		return resp, err
	}

	err = s.profilesBusiness.CheckPermission(ctx, profileID, characterID, metricID)
	if err != nil {
		return resp, err
	}

	resp.Authorized = true

	return resp, nil
}
