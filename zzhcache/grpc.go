package zzhcache

import (
	"context"
	"fmt"
	"goog
	"log"
	"net"
	"sync"
	"zzhcache/consistenthash"
	pb "zzhcache/zzhcachepb"
)

type GRPCPool struct {
	self        string
	mu          sync.Mutex
	peers       *consistenthash.Map
	grpcGetters map[string]*grpcGetter
}

func NewGRPCPool(self string) *GRPCPool {
	return &GRPCPool{
		self:        self,
		grpcGetters: make(map[string]*grpcGetter),
	}
}

func (p *GRPCPool) Log(format string, v ...interface{}) {
	log.Printf("[GRPC Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *GRPCPool) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterGroupCacheServer(server, &grpcServer{pool: p})

	p.Log("GRPC server started at %s", addr)
	return server.Serve(lis)
}

func (p *GRPCPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.grpcGetters = make(map[string]*grpcGetter, len(peers))
	for _, peer := range peers {
		p.grpcGetters[peer] = &grpcGetter{addr: peer}
	}
}

func (p *GRPCPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.grpcGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*GRPCPool)(nil)

type grpcServer struct {
	pb.UnimplementedGroupCacheServer
	pool *GRPCPool
}

func (s *grpcServer) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	group := GetGroup(req.GetGroup())
	if group == nil {
		return nil, fmt.Errorf("group %s not found", req.GetGroup())
	}

	view, err := group.Get(req.GetKey())
	if err != nil {
		return nil, err
	}

	return &pb.Response{Value: view.ByteSlice()}, nil
}

type grpcGetter struct {
	addr string
}

func (g *grpcGetter) Get(in *pb.Request, out *pb.Response) error {
	conn, err := grpc.Dial(g.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewGroupCacheClient(conn)
	resp, err := client.Get(context.Background(), in)
	if err != nil {
		return err
	}

	out.Value = resp.Value
	return nil
}

var _ PeerGetter = (*grpcGetter)(nil)
