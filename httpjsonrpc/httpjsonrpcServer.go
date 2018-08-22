package httpjsonrpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ServerMux struct {
	m               map[string]func(*http.Request, map[string]interface{}) map[string]interface{}
	defaultFunction func(http.ResponseWriter, *http.Request)
}

var mainMux ServerMux

func InitServeMux() {
	mainMux = make(map[string]func(*http.Request, map[string]interface{}) map[string]interface{})
}

func HandleFunc(pattern string, handler func(*http.Request, map[string]interface{}) map[string]interface{}) {
	mainMux.m[pattern] = handler
}

func SetDefaultFunc(def func(http.ResponseWriter, *http.Request)) {
	mainMux.defaultFunction = def
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		if mainMux.defaultFunction != nil {
			log.Printf("HTTP JSON RPC Handle - Method!=\"POST\"")
			mainMux.defaultFunction(w, r)
			return
		} else {
			log.Panicf("HTTP JSON RPC Handle - Method!=\"POST\"")
			return
		}
	}

	if r.Body == nil {
		if mainMux.defaultFunction != nil {
			log.Printf("HTTP JSON RPC Handle - Request body is nil")
			mainMux.defaultFunction(w, r)
			return
		} else {
			log.Panicf("HTTP JSON RPC Handle - Request body is nil")
			return
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("HTTP JSON RPC Handle - ioutil.ReadAll: %v", err)
		return
	}

	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatalf("HTTP JSON RPC Handle - json.Unmarshal: %v", err)
		return
	}
	function, ok := mainMux.m[request["method"].(string)]
	if ok {
		response := function(r, request)
		data, err := json.Marshal(response)
		if err != nil {
			log.Fatalf("HTTP JSON RPC Handle - json.Marshal:%v", err)
			return
		}
		w.Write(data)
	} else {
		log.Println("HTTP JSON RPC Handle - No function to call for", request["method"])
		data, err := json.Marshal(map[string]interface{}{
			"result": nil,
			"error": map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
				"data":    "The called method was not found on the server",
			},
			"id": request["id"],
		})
		if err != nil {
			log.Fatalf("HTTP JSON RPC Handle - json.Marshal: %v", err)
			return
		}
		w.Write(data)
	}
}
