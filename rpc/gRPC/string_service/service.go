package string_service

import (
	"context"
	"rpc/gRPC/pb"
	"strings"
)

type StringService struct {
}

func (s StringService) Concat(ctx context.Context, request *pb.StringRequest) (*pb.StringResponse, error) {
	if len(request.A)+len(request.B) > 1024 {
		response := pb.StringResponse{Ret: ""}
		return &response, nil
	}
	response := pb.StringResponse{Ret: request.A + request.B}
	return &response, nil
}

func (s StringService) Diff(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	if len(req.A) < 1 || len(req.B) < 1 {
		return &pb.StringResponse{Ret: ""}, nil
	}

	res := ""
	if len(req.A) >= len(req.B) {
		for _, char := range req.B {
			if strings.Contains(req.A, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range req.A {
			if strings.Contains(req.B, string(char)) {
				res = res + string(char)
			}
		}
	}
	return &pb.StringResponse{Ret: res}, nil
}
