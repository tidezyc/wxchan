package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/tidezyc/wxchan/weixin"
)

func main() {
	cli := weixin.NewweixinClient()
	err := cli.Login()
	if err != nil {
		log.Fatalf("weixin login err:%s", err)
	}
	contacts, err := cli.GetContacts()
	if err != nil {
		log.Fatalf("weixin get contacts err:%s", err)
	}
	for _, contact := range contacts {
		fmt.Println(contact)
	}
	http.HandleFunc("/send", send)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("", "9086"), nil))
}

func send(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
