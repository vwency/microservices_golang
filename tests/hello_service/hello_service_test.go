package hello_service_test

import (
	"context"
	"testing"
	"net"

	"github.com/stretchr/testify/assert"
	"github.com/vwency/microservices_golang/proto/hello_service"
	"github.com/vwency/microservices_golang/internal/hello_service/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func newTestServer() *grpc.Server {
	server := grpc.NewServer()
	helloHandler := handler.NewHelloHandler()
	hello_service.RegisterHelloServiceServer(server, helloHandler)
	return server
}

func TestSayHello(t *testing.T) {
	lis := bufconn.Listen(bufSize)

	server := newTestServer()
	go func() {
		if err := server.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := hello_service.NewHelloServiceClient(conn)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"test with Hello", "Hello", "hello"},
		{"test with hello", "hello", "hello"},
		{"test with HELLO", "HELLO", "hello"},
		{"test with mixed case", "HeLlO", "hello"},
		{"test with no hello", "None", "None"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.SayHello(context.Background(), &hello_service.HelloRequest{
				Text: tt.input,
			})
			if err != nil {
				t.Fatalf("could not greet: %v", err)
			}

			assert.Equal(t, tt.expected, resp.GetText())
		})
	}
}
