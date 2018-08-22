package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Arith int

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Result int

type GetBestBlockHashResp struct {
	Id     interface{} `json:"id"`
	Result Reply       `json:"result"`
	Error  interface{} `json:"error"`
}

func (t *Arith) Multiply(args *Args, result *Result) error {
	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	print("test the multipy\n")
	*result = Result(args.A * args.B)
	return nil
}

func (t *Arith) getbestblockhash(args *Args, result *GetBestBlockHashResp) error {
	log.Printf("Multiply with \n")
	print("test")
	return nil
}

func startServer() {
	server := rpc.NewServer()
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	l, e := net.Listen("tcp", ":10333")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	go startServer()
	conn, err := net.Dial("tcp", "localhost:10333")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	args := &Args{7, 8}
	var reply Reply
	c := jsonrpc.NewClient(conn)
	dec := json.NewDecoder(conn)

	for i := 0; i < 1; i++ {
		err := c.Call("Arith.Multiply", args, &reply)
		fmt.Fprint(conn, `{"jsonrpc": "2.0", "method": "getbestblockhash", "params": [], "id": 2`)
		var resp GetBestBlockHashResp
		err := dec.Decode(&resp)
		if err != nil {
			log.Fatal("Decode:%s:", err)
		}
		fmt.Printf("Get best block hash resp:%d %d\n", resp.Id, resp.Result)
		print("test the end!\n")
	}
}
