package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. Load the CA certificate used to sign the client certificates
	clientCACert, err := os.ReadFile("./resources/client_ca.crt") // CA that signed client certs
	if err != nil {
		log.Fatalf("Failed to read client CA: %v", err)
	}
	clientCAPool := x509.NewCertPool()
	clientCAPool.AppendCertsFromPEM(clientCACert)

	// 2. Load the server certificate and key
	serverCert, err := tls.LoadX509KeyPair("./resources/server.crt", "./resources/server.key")
	if err != nil {
		log.Fatalf("Failed to load server key pair: %v", err)
	}

	// 3. Configure the TLS connection
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // This is the crucial mTLS setting
		ClientCAs:    clientCAPool,                   // Trust client certs signed by this CA
		MinVersion:   tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler:   http.HandlerFunc(handler),
	}

	log.Printf("Starting server on https://localhost:8443")
	// Use ListenAndServeTLS with empty strings for certFile/keyFile since they're in TLSConfig
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Optional: Inspect the client certificate used by the current request
	if len(r.TLS.PeerCertificates) > 0 {
		cert := r.TLS.PeerCertificates[0]
		log.Printf("Request from Subject: %s", cert.Subject.CommonName)
	}
	w.Write([]byte("Secure response!"))
}
