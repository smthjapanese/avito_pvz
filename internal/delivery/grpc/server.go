package grpc

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc"

	pbv1 "github.com/smthjapanese/avito_pvz/github.com/avito_pvz/pvz/pvz_v1"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
)

type Server struct {
	pbv1.UnimplementedPVZServiceServer
	pvzUseCase usecase.PVZUseCase
	grpcServer *grpc.Server
	wg         sync.WaitGroup
}

func NewServer(pvzUseCase usecase.PVZUseCase) *Server {
	return &Server{
		pvzUseCase: pvzUseCase,
		grpcServer: grpc.NewServer(),
	}
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

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	pbv1.RegisterPVZServiceServer(s.grpcServer, s)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("gRPC server listening on %s", port)
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Printf("Failed to serve: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
		s.wg.Wait()
	}
}

func Start(pvzUseCase usecase.PVZUseCase, port string) {
	server := NewServer(pvzUseCase)
	if err := server.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
