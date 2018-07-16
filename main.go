/**/
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
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
		Allowed bool            `json:"allowed"`
		Status  AdmissionStatus `json:"status"`
	} `json:"response"`
}

// AdmissionStatus is baz
type AdmissionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Code    int    `json:"code"`
}

// AdmissionResponse is foo
type AdmissionResponse struct {
	Kind    string `json:"kind"`
	Request struct {
		Object struct {
			Spec struct {
				Containers []struct {
					Name string
					Env  []struct {
						Name  string
						Value string
					}
				}
			}
		}
	}
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info("validating")

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	admissionStatus := new(AdmissionStatus)
	admissionReview := &AdmissionReview{
		Response: struct {
			Allowed bool            `json:"allowed"`
			Status  AdmissionStatus `json:"status"`
		}{
			Allowed: true,
			Status:  *admissionStatus,
		},
	}

	// var admissionRequest = req.body
	// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
	// var body []byte
	// if r.Body != nil {
	// 	if data, err := ioutil.ReadAll(r.Body); err == nil {
	// 		body = data
	// 	}
	// }
	// fmt.Println("")
	// fmt.Println("", body)

	fmt.Println("body:")
	admissionResponse := new(AdmissionResponse)
	json.NewDecoder(r.Body).Decode(admissionResponse)
	fmt.Println(admissionResponse)
	foundEnv := false
	for _, container := range admissionResponse.Request.Object.Spec.Containers {
		// fmt.Println("index:", index, " ", "len:", len(container.Env))

		if len(container.Env) > 0 {
			foundEnv = true
			admissionStatus.Status = "Failure"
			admissionStatus.Message = "The container \"" + container.Name + "\" is using env vars"
			admissionStatus.Reason = "The container \"" + container.Name + "\" is using env vars"
			admissionStatus.Code = 402
			break
		}
	}

	if foundEnv {
		admissionReview.Response.Allowed = false
		admissionReview.Response.Status = *admissionStatus
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
