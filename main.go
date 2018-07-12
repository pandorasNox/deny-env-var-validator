/**/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"log"
	"crypto/tls"

	"github.com/golang/glog"
)

// type AdmissionStatus struct {
// 	status: string
// 	message: string
// 	reason: string
// 	code: int
// }

type AdmissionResponse struct {
	Allowed	bool `json:"allowed"`
	// status	AdmissionStatus
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info("validating")

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	admissionResponse := AdmissionResponse{Allowed:true}

	// js, err := json.Marshal(admissionResponse)
	// if err != nil {
	//   http.Error(w, err.Error(), http.StatusInternalServerError)
	//   return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admissionResponse)
}

func main() {
	// var config Config
	// config.addFlags()
	// flag.Parse()

	http.HandleFunc("/content", serveContent)

	// // clientset := getClient()
	// server := &http.Server{
	// 	// Addr:      ":443",
	// 	Addr:      ":8083",
	// 	// TLSConfig: configTLS(config, clientset),
	// }
	// // server.ListenAndServe()

	// ==============
	// use TLS
	// ==============

	cfg := &tls.Config{
		MinVersion:					tls.VersionTLS12,
		CurvePreferences:			[]tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites:	true,
		CipherSuites: []uint16{
			// tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// clientset := getClient()
	server := &http.Server{
		// Addr:      ":443",
		Addr:      ":8083",
		// TLSConfig: configTLS(config, clientset),
		TLSConfig: cfg,
	}

	err := server.ListenAndServeTLS("/certs/ssl-cert.pem", "/certs/ssl-key.pem")
	log.Fatal(err)
}
