package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
)

func main() {
	caCert, err := ioutil.ReadFile("/var/app/server.crt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CA DATA", string(caCert))
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cfg := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}
	http.HandleFunc("/index.yaml", handler)
	srv := &http.Server{
		Addr:      ":8443",
		TLSConfig: cfg,
	}
	log.Fatal(srv.ListenAndServeTLS("/var/app/server.crt", "/var/app/server.key"))
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://azure-samples.github.io/helm-charts/index.yaml")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

