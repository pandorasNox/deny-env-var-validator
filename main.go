/**/
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/golang/glog"
)

// type AdmissionStatus struct {
// 	status: string
// 	message: string
// 	reason: string
// 	code: int
// }

// AdmissionReview returns a validation to kubernetes api server
type AdmissionReview struct {
	Response struct {
		Allowed bool `json:"allowed"`
		// status	AdmissionStatus
	} `json:"response"`
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info("validating")

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	// var admissionRequest = req.body
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	fmt.Println([]byte(body))

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	admissionReview := &AdmissionReview{
		Response: struct {
			Allowed bool `json:"allowed"`
		}{
			Allowed: true,
		},
	}

	// js, err := json.Marshal(admissionReview)
	// if err != nil {
	//   http.Error(w, err.Error(), http.StatusInternalServerError)
	//   return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admissionReview)
}

func main() {
	var tlsDisabled *bool
	tlsDisabled = flag.Bool("tlsDisabled", false, "(optional) disables tls for the server")
	flag.Parse()

	// var config Config
	// config.addFlags()
	// flag.Parse()

	fmt.Println("tlsDisabled: ", bool(*tlsDisabled))

	http.HandleFunc("/content", serveContent)

	if bool(*tlsDisabled) {
		server := &http.Server{
			Addr: ":8083",
		}
		server.ListenAndServe()
		os.Exit(0)
	}

	// ==============
	// use TLS
	// ==============

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
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
		Addr: ":8083",
		// TLSConfig: configTLS(config, clientset),
		TLSConfig: cfg,
	}

	err := server.ListenAndServeTLS("/certs/ssl-cert.pem", "/certs/ssl-key.pem")
	log.Fatal(err)
}
