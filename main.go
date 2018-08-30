package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ShardingBlockchain/httpjsonrpc"
)

func getInfo(req *http.Request, cmd map[string]interface{}) map[string]interface{} {
	id := cmd["id"].(float64)
	fmt.Println(id)

	params := cmd["params"]

	fmt.Println(params)

	if cmd == nil {
		cmd = make(map[string]interface{})
	}

	return cmd
}

func init() {
	go start_server()
}

func sendToAddress(req *http.Request, cmd map[string]interface{}) map[string]interface{} {
	return cmd
}

func start_server() {
	print("hello start server\n")
	httpjsonrpc.InitServeMux()
	http.HandleFunc("/", httpjsonrpc.Handle)
	httpjsonrpc.HandleFunc("getinfo", getInfo)
	httpjsonrpc.HandleFunc("sendtoaddress", sendToAddress)
	err := http.ListenAndServe("localhost:10332", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())

	}
}

func start_client() {
	var res map[string]interface{}
	var err error
	print("start client\n")

	res, err = httpjsonrpc.Call("http://127.0.0.1:10332", "getinfo", 1, []interface{}{})
	if err != nil {
		log.Fatalf("Err:%v", err)
	}
	log.Println(res)
	/*
		// call send to address
		params := []interface{}{"asset_id", "address", 56}
		res, err = httpjsonrpc.Call("http://127.0.0.1:10332", "sendtoaddress", 2, params)
		if err != nil {
			log.Fatalf("Err: %v", err)
		}
		log.Println(res)
	*/
}

func main() {
	time.Sleep(2 * time.Second)
	start_client()
}
