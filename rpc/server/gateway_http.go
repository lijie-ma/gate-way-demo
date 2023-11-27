package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	helloworldpb "gate-way-demo/proto/helloworld"
)

type serverHttp struct {
	helloworldpb.UnimplementedGreeterServer
}

func NewServerHttp() *serverHttp {
	return &serverHttp{}
}

func (s *serverHttp) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("md", data.Get("auth"), len(data.Get("auth")))
	}
	return &helloworldpb.HelloReply{Message: in.Name + " world"}, nil
}

func runHttp() http.Handler {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	mux := runtime.NewServeMux(
		// convert header in response(going from gateway) from metadata received.
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			log.Println("header", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			//creating a new HTTTPStatusError with a custom status, and passing error
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// using default handler to do the rest of heavy lifting of marshaling error and adding headers
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}))
	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	helloworldpb.RegisterGreeterServer(s, NewServerHttp())
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// Register Greeter
	err = helloworldpb.RegisterGreeterHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	return mux
}

func RunHttpWrapGin() {
	mux := runHttp()
	httpServer := gin.New()
	httpServer.Use(gin.Logger())
	httpServer.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(Auth(mux)))
	httpServer.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: httpServer,
		// Handler: MultipleMiddleware(mux, Auth())
	}

	log.Println("gin gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

func RunWithHttp() {
	// Create a listener on TCP port
	mux := runHttp()
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: Auth(mux),
		// Handler: MultipleMiddleware(mux, Auth())
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
