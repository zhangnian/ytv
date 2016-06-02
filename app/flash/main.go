package main

import (
	"io"
	"log"
	"net/http"
)

func CrossDomain(w http.ResponseWriter, r *http.Request) {
	hp := `<?xml version="1.0" ?>
           <!DOCTYPE cross-domain-policy SYSTEM "http://www.adobe.com/xml/dtds/cross-domain-policy.dtd">
            <cross-domain-policy>
            <site-control permitted-cross-domain-policies="all"/>
            <allow-access-from domain="*.pztai.cn"/>
            <allow-http-request-headers-from domain="*.pztai.cn" headers="*"/>
            </cross-domain-policy>`
	io.WriteString(w, hp)
}

func main() {
	http.HandleFunc("/crossdomain.xml", CrossDomain)
	err := http.ListenAndServe(":843", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
