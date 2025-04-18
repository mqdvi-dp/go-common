package rpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/factory"
	log "github.com/mqdvi-dp/go-common/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type rpcServer struct {
	opt          option
	serverEngine *grpc.Server
	listener     net.Listener
	service      factory.ServiceFactory
}

// New create a new gRPC server
func New(service factory.ServiceFactory, opts ...OptionFunc) factory.AppServerFactory {
	var (
		keepAliveEnforce = keepalive.EnforcementPolicy{
			MinTime:             env.GetDuration("GRPC_MIN_TIME", time.Duration(10)*time.Second),
			PermitWithoutStream: true,
		}
		keepAliveServer = keepalive.ServerParameters{
			MaxConnectionIdle: env.GetDuration(
				"GRPC_MAX_CONNECTION_IDLE_DURATION",
				time.Duration(10)*time.Second,
			), // if a client is idle for 10s, send a go away
			MaxConnectionAge: env.GetDuration(
				"GRPC_MAX_CONNECTION_AGE",
				time.Duration(30)*time.Second,
			), // allows 30s for pending RPCs to complete before forcibly closing connections
			MaxConnectionAgeGrace: env.GetDuration(
				"GRPC_MAX_CONNECTION_AGE_GRACE",
				time.Duration(10)*time.Second,
			), // allows 10s for pending RPCs to complete before forcibly closing connections
			Time: env.GetDuration(
				"GRPC_TIME_PING_CLIENT",
				time.Duration(5)*time.Second,
			), // ping the client if it is idle for 5s to ensure the connection is still alive
			Timeout: env.GetDuration(
				"GRPC_TIMEOUT",
				time.Duration(1)*time.Second,
			), // wait 1s for the ping ack before assuming the connection is dead
		}
	)
	intercept := &interceptor{serviceName: service.Name()}

	server := &rpcServer{
		service: service,
		opt:     getDefaultOption(),
		serverEngine: grpc.NewServer(
			grpc.KeepaliveEnforcementPolicy(keepAliveEnforce),
			grpc.KeepaliveParams(keepAliveServer),
			grpc.UnaryInterceptor(
				intercept.chainUnaryServer(
					intercept.unaryServerTracerInterceptor,
				),
			),
		),
	}
	reflection.Register(server.serverEngine)

	for _, opt := range opts {
		opt(&server.opt)
	}

	port := server.opt.tcpPort
	var err error
	server.listener, err = net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	intercept.opt = &server.opt
	if h := service.GRPCHandler(); h != nil {
		h.RegisterGrpcServer(server.serverEngine)
	}

	for root, info := range server.serverEngine.GetServiceInfo() {
		for _, method := range info.Methods {
			log.Blue(fmt.Sprintf(`[GRPC-METHOD] (root): %-10s (method): %-8s (metadata): %v`, `"`+root+`"`, `"`+method.Name+`"`, info.Metadata))
		}
	}

	return server
}

func (r *rpcServer) Serve() {
	log.BlueBold(fmt.Sprintf("⇨ GRPC Server run at port [::]%s", r.opt.tcpPort))

	if err := r.serverEngine.Serve(r.listener); err != nil {
		log.Log.Fatal(err)
	}
}

func (r *rpcServer) Shutdown(_ context.Context) {
	defer fmt.Printf("\x1b[33;1m%s Stopping GRPC server:\x1b[0m \x1b[32;1mSUCCESS\x1b[0m\n", r.service.Name())

	r.serverEngine.GracefulStop()
	_ = r.listener.Close()
}

func (r *rpcServer) Name() string {
	return string(constants.GRPC)
}
