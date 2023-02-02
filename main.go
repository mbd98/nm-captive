package main

import (
	"github.com/Wifx/gonetworkmanager/v2"
	"log"
	"net/http"
)

func determineUrl() {
	req, err := http.NewRequest(http.MethodHead, "http://detectportal.firefox.com/canonical.html", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0")
	if err != nil {
		log.Panicln(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(resp.Header)
}

func main() {
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		log.Fatalln(err)
	}
	defer nm.Unsubscribe()
	ch := nm.Subscribe()
	for {
		sig := <-ch
		if sig.Name == "org.freedesktop.NetworkManager.Device.StateChanged" {
			err := nm.CheckConnectivity()
			if err != nil {
				log.Panicln(err)
			}
			connectivity, err := nm.GetPropertyConnectivity()
			if err != nil {
				log.Panicln(err)
			}
			if connectivity == gonetworkmanager.NmConnectivityPortal {
				log.Println("Captive portal detected")
				go determineUrl()
			}
		}
	}
}
