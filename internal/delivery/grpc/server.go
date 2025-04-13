package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"

	"google.golang.org/grpc"

	pbv1 "github.com/smthjapanese/avito_pvz/github.com/avito_pvz/pvz/pvz_v1"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
)

type Server struct {
	pbv1.UnimplementedPVZServiceServer
	pvzUseCase usecase.PVZUseCase
}

func NewServer(pvzUseCase usecase.PVZUseCase) *Server {
	return &Server{pvzUseCase: pvzUseCase}
}

// GetPVZList реализует gRPC метод для получения списка ПВЗ
func (s *Server) GetPVZList(ctx context.Context, req *pbv1.GetPVZListRequest) (*pbv1.GetPVZListResponse, error) {
	pvzs, err := s.pvzUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := &pbv1.GetPVZListResponse{}
	for _, pvz := range pvzs {
		response.Pvzs = append(response.Pvzs, &pbv1.PVZ{
			Id:               pvz.ID.String(),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             string(pvz.City),
		})
	}

	return response, nil
}

func Start(pvzUseCase usecase.PVZUseCase, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pbv1.RegisterPVZServiceServer(grpcServer, NewServer(pvzUseCase))

	log.Printf("gRPC server listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
