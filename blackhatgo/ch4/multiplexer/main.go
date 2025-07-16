package main

import (
	"net/http/httputil"
)

var (
	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy)
)

func init() {

}

func main() {

}
