package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"	// https://pkg.go.dev/github.com/godbus/dbus/v5
	"log"
)

const (
	// https://networkmanager.dev/docs/api/latest/gdbus-org.freedesktop.NetworkManager.html#gdbus-property-org-freedesktop-NetworkManager.Connectivity
	NM_CONNECTIVITY = "Connectivity"
	
	// https://networkmanager.dev/docs/api/latest/nm-dbus-types.html#NMConnectivityState
	NM_CONNECTIVITY_PORTAL = 2
	NM_OBJPATH = "/org/freedesktop/NetworkManager"
	
	// https://dbus.freedesktop.org/doc/dbus-specification.html#standard-interfaces-properties
	DBUS_IFACE_PROP = "org.freedesktop.DBus.Properties"
	DBUS_SIG_PROPCHANGE = "org.freedesktop.DBus.Properties.PropertiesChanged"
)

func main() {
	log.Println("Connecting to system bus...")
	bus, err := dbus.SystemBus()
	if err != nil {
		log.Fatalf("Failed to connect to system bus: %v\n", err)
	}
	defer bus.Close()
	log.Println("Setup PropertiesChanged signal watch...")
	if err = bus.AddMatchSignal(dbus.WithMatchObjectPath(NM_OBJPATH), dbus.WithMatchInterface(DBUS_IFACE_PROP)); err != nil {
		log.Panicf("Failed to listen for signal: %v\n", err)
	}
	ch := make(chan *dbus.Signal, 16)
	bus.Signal(ch)
	log.Println("Listening for signals...")
	for sig := range ch {
		if sig.Name == DBUS_SIG_PROPCHANGE {
			m := sig.Body[1].(map[string]dbus.Variant)
			if connvar, ok := m[NM_CONNECTIVITY]; ok && connvar.Value().(uint32) == NM_CONNECTIVITY_PORTAL {
				fmt.Println("Portal detected.  Open a web browser to log in")
			}
		}
	}
}
