package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"log"
)

const (
	NM_CONNECTIVITY = "Connectivity"
	NM_CONNECTIVITY_PORTAL = 2
)

func main() {
	log.Println("Connecting to system bus...")
	bus, err := dbus.SystemBus()
	if err != nil {
		log.Fatalf("Failed to connect to system bus: %v\n", err)
	}
	defer bus.Close()
	log.Println("Setup PropertiesChanged signal watch...")
	if err = bus.AddMatchSignal(dbus.WithMatchObjectPath("/org/freedesktop/NetworkManager"), dbus.WithMatchInterface("org.freedesktop.DBus.Properties")); err != nil {
		log.Panicf("Failed to listen for signal: %v\n", err)
	}
	ch := make(chan *dbus.Signal, 16)
	bus.Signal(ch)
	log.Println("Listening for signals...")
	for sig := range ch {
		if sig.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			m := sig.Body[1].(map[string]dbus.Variant)
			if connvar, ok := m[NM_CONNECTIVITY]; ok && connvar.Value().(uint32) == NM_CONNECTIVITY_PORTAL {
				fmt.Println("Portal detected.  Open a web browser to log in")
			}
		}
	}
}
