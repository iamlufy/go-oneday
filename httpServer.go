package main

import (
	"go-oneday/odserver"
	"net/http"
	"fmt"
	"net/url"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	var params url.Values = req.URL.Query()
	fmt.Println(params)
	fmt.Fprint(w, "hello world")
}
func HelloServer2(c *odserver.Context) {

	fmt.Fprint(c.Rw, "hello world test2")
}


func HelloServer3(c *odserver.Context) {

	fmt.Fprint(c.Rw, c.Params)
}
func HelloServer4(c *odserver.Context) {

	fmt.Fprint(c.Rw, "hello world HelloServer4")
}


func main() {
	o := odserver.Default()
	o.Start("/main").
		Target("/test/").Get(HelloServer).Post(HelloServer).Delete(HelloServer).And().
		Target("/test2").Get(HelloServer2)
	o.Start("/{test}/main/").Target("/number/{number}").
		Get(HelloServer3).Post(HelloServer4)

	http.ListenAndServe(":8080",o)

}
