package server

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func UserRpc(w http.ResponseWriter, r *http.Request) {
	rpc.RegisterName("HelloService", new(HelloService))
	var conn io.ReadWriteCloser = struct {
		io.Writer
		io.ReadCloser
	}{
		ReadCloser: r.Body,
		Writer:     w,
	}

	rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
}
