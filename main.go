/***************************************************************
 *	Created by: Rahul Shewale
 *	Date: 20/12/2018
 ***************************************************************/
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! Letsencrypt\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

//Examples with Auto generate SSL Certificate
func main2() {
	//API Router
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	//Set your domain name and certificate directory path
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("pixel.gentlereader.com"), //Your domain here
		Cache:      autocert.DirCache("/home/pixel/certs"),           //Folder for storing certificates
	}
	//443 = https and 80 = http
	server := &http.Server{
		Addr:    ":443",
		Handler: router,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	//this Listen server for redirect https
	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

	log.Fatal(server.ListenAndServeTLS("", ""))
}

//Example of manually generated server.crt and private.key for domain
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://pixel.gentlereader.com:443"+r.RequestURI, http.StatusMovedPermanently)
}
func main() {
	//API Router
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	//443 = https and 80 = http
	server := &http.Server{
		Addr:    ":443",
		Handler: router,
	}
	//this Listen server for redirect https
	go server.ListenAndServeTLS("/home/pixel/certs_client/server.crt", "/home/pixel/certs_client/private.key")

	if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}

}
