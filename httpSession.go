package main

import (
	"net/http"
	"fmt"
	"./odsession"
	"log"
)

var i int64 = 0


func readGoSession(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie(cookieName)
	//cookie失效
	if err != nil {
		fmt.Fprintf(w, "nil")
		return
	}
	s := manager.Session(c.Value)
	fmt.Print(s.Load("test"))
	fmt.Fprintf(w, c.Value)
}
func expireGoSession(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie(cookieName)
	s := manager.Session(c.Value)
	manager.RemoveSession(s.SessionId(), w)

}

var manager odsession.Manager

const cookieName = "gosessionid"
const maxLifeTime = 1800

func HelloServer1(w http.ResponseWriter, req *http.Request) {

	fmt.Fprint(w,"hello world")
}


func main() {
	p := odsession.NewMemoryProvider()

	manager, _ = odsession.NewManager(p, cookieName, maxLifeTime) //半个钟

	http.HandleFunc("/test", HelloServer)
	http.HandleFunc("/read", readGoSession)
	http.HandleFunc("/expire", expireGoSession)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
