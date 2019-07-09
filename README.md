#Example of Create API with https (Let’s Encrypt)

#Requirements
1) Virtual Machine on Digital Ocean
2) Domain name and point to your server (VM)

#Let’s Encrypt and the ACME protocol
Let’s Encrypt is a very well known and trusted SSL certificate issuer that provides a free and automated generation process. It is possible to issue a certificate in less than a second without any registration process or payment.

Autocert is a Go package that implements the ACME protocol used to generates certificates on Let’s Encrypt. This is the only package dependency that you will need, no other installation or package is required.

You can get it as any other Go package.

```
go get golang.org/x/crypto/acme/autocert
```

#Examples with Auto generate SSL Certificate

```go
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

func main() {
	//API Router
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	//Set your domain name and certificate directory path
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("sub.domain.com"),       //Your domain or subdomain here
		Cache:      autocert.DirCache("/home/projectfolder/certs"), //Folder for storing certificates
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
```



#Example of manually generated server.crt and private.key for domain 
```go
package main
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://domain.com:443"+r.RequestURI, http.StatusMovedPermanently)
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
	go server.ListenAndServeTLS("/home/projectname/certs_client/server.crt", "/home/projectname/certs_client/private.key")

	if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}

}
```

#Reference Link
1) https://goenning.net/2017/11/08/free-and-automated-ssl-certificates-with-go/
2) https://godoc.org/golang.org/x/crypto/acme/autocert
3 https://www.digitalocean.com/community/tutorials/how-to-secure-apache-with-let-s-encrypt-on-ubuntu-14-04 (Note: This link only for make https url using apache server)
4)https://stackoverflow.com/questions/37321760/how-to-set-up-lets-encrypt-for-a-go-server-application

#Generate SSL Certificate for domain
http://www.selfsignedcertificate.com/