// Code generated by protoc-gen-go-client. DO NOT EDIT.
// Sources: queue.proto

package client

import (
	context "context"

	grpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	pb "github.com/erda-project/erda-proto-go/core/pipeline/queue/pb"
	grpc1 "google.golang.org/grpc"
)

// Client provide all service clients.
type Client interface {
	// QueueService queue.proto
	QueueService() pb.QueueServiceClient
}

// New create client
func New(cc grpc.ClientConnInterface) Client {
	return &serviceClients{
		queueService: pb.NewQueueServiceClient(cc),
	}
}

type serviceClients struct {
	queueService pb.QueueServiceClient
}

func (c *serviceClients) QueueService() pb.QueueServiceClient {
	return c.queueService
}

type queueServiceWrapper struct {
	client pb.QueueServiceClient
	opts   []grpc1.CallOption
}

func (s *queueServiceWrapper) CreateQueue(ctx context.Context, req *pb.QueueCreateRequest) (*pb.QueueCreateResponse, error) {
	return s.client.CreateQueue(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *queueServiceWrapper) GetQueue(ctx context.Context, req *pb.QueueGetRequest) (*pb.QueueGetResponse, error) {
	return s.client.GetQueue(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *queueServiceWrapper) PagingQueue(ctx context.Context, req *pb.QueuePagingRequest) (*pb.QueuePagingResponse, error) {
	return s.client.PagingQueue(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *queueServiceWrapper) UpdateQueue(ctx context.Context, req *pb.QueueUpdateRequest) (*pb.QueueUpdateResponse, error) {
	return s.client.UpdateQueue(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}

func (s *queueServiceWrapper) DeleteQueue(ctx context.Context, req *pb.QueueDeleteRequest) (*pb.QueueDeleteResponse, error) {
	return s.client.DeleteQueue(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}