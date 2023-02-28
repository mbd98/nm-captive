package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"log"
)

func main() {
	bus, err := dbus.SystemBus()
	if err != nil {
		log.Fatalf("Failed to connect to system bus: %v\n", err)
	}
	defer bus.Close()
	if err = bus.AddMatchSignal(dbus.WithMatchObjectPath("/org/freedesktop/NetworkManager"), dbus.WithMatchInterface("org.freedesktop.DBus.Properties")); err != nil {
		log.Fatalf("Failed to listen for signal: %v\n", err)
	}
	ch := make(chan *dbus.Signal, 16)
	bus.Signal(ch)
	for sig := range ch {
		if sig.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			fmt.Println(sig.Body)
		}
	}
}
