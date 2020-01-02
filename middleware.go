package main

import (
	"fmt"
	"github.com/cloudflare/cfrpki/validator/lib"
	"io/ioutil"
	"log"
	"net/http"

	"bytes"
	"crypto/tls"

	"flag"
)

var (
	bind               = flag.String("bind", ":3002", "Address to listen")
	bindEnableHTTPs    = flag.Bool("bind.https", false, "Enable https")
	bindTLSCertificate = flag.String("bind.tls.cert", "server.pem", "TLS server certificate")
	bindTLSKey         = flag.String("bind.tls.key", "server.key", "TLS server key")

	proxy            = flag.String("proxy", "127.0.0.1:3001", "Address of the proxy")
	proxyEnableHTTPs = flag.Bool("proxy.https", true, "Enable https for proxy")
	proxyTLSVerify   = flag.Bool("proxy.verify", false, "Enable TLS validation for proxy")
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("-> %v %v\n", r.Method, r.URL.String())

	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		data, err := librpki.DecodeXML(body)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("%v\n", string(data.Content))
		}

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: !*proxyTLSVerify,
				},
			},
		}

		bodyOutBuf := bytes.NewBuffer(body)

		dstUri := r.URL
		if *proxyEnableHTTPs {
			dstUri.Scheme = "https"
		} else {
			dstUri.Scheme = "http"
		}

		dstUri.Host = *proxy

		resp, err := client.Post(dstUri.String(), "application/rpki-publication", bodyOutBuf)
		fmt.Printf("<- POST %v %v\n", dstUri.String(), resp.Status)

		if err != nil {
			fmt.Printf("%v\n", err)
		}

		if resp.Body != nil {
			body, _ = ioutil.ReadAll(resp.Body)
			data, err = librpki.DecodeXML(body)
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				fmt.Printf("%v\n", string(data.Content))
			}
			w.Write(body)
		}
	}
}

func main() {
	flag.Parse()
	log.Printf("Serving on %v\n", *bind)

	http.HandleFunc("/", handler)

	if *bindEnableHTTPs {
		log.Fatal(http.ListenAndServeTLS(*bind, *bindTLSCertificate, *bindTLSKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(*bind, nil))
	}

}
