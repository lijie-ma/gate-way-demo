package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	pb "gate-way-demo/proto/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var duration time.Duration = 3 * time.Second

func TestRpc(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), duration)
	conn, err := GetConn(ctx)
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	t.Log(client.SayHello(ctx, &pb.HelloRequest{Name: "client hello"}))
}

func getConnAddr() string {
	return "localhost:8080"
}

func GetConn(ctx context.Context) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithCredentialsBundle(insecure.NewBundle()))
	return grpc.DialContext(ctx, getConnAddr(), opts...)
}

func testHttp(t *testing.T, headers http.Header) {
	data := []byte(`{"name":"http hello"}`)
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8090/v1/example/echo", bytes.NewBuffer(data))
	if nil != err {
		t.Errorf("http req err: %+v", err)
		return
	}
	req.Header = headers
	client := http.Client{}
	resp, err := client.Do(req)
	if nil != err {
		t.Errorf("http client err: %+v", err)
		return
	}
	var buf bytes.Buffer
	var respBody = make([]byte, 5)
	for {
		n, err := resp.Body.Read(respBody)
		if nil != err {
			if errors.Is(err, io.EOF) {
				buf.WriteString(string(respBody[:n]))
				break
			}
			t.Errorf("http read body err: %+v", err)
			break
		}
		buf.WriteString(string(respBody))
	}

	defer resp.Body.Close()
	if nil != err {
		t.Errorf("http read body err: %+v", err)
		return
	}
	t.Logf("resp: %s", buf.String())
}

func TestHttp(t *testing.T) {
	unauth := make(http.Header)
	unauth.Add("content-type", "application/json")
	testHttp(t, unauth)
	t.Log("----------------")
	auth := make(http.Header)
	auth.Add("content-type", "application/json")
	auth.Add("Authorization", "bearer test_auth")
	testHttp(t, auth)
}
