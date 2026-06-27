package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	_ "embed"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus/v5"
	"github.com/skratchdot/open-golang/open"
)

const (
	GAMEISREADY     = "Your game is ready"
	GAMEISUNPAUSING = "The game is unpausing..."
	dota2ID         = "Dota2"
	homePage        = "https://github.com/dastier/GoDota2Accepter"
)

//go:embed assets/d_enabled.svg
var iconEnabled []byte

//go:embed assets/d_disabled.svg
var iconDisabled []byte

var (
	enabled        uint32 = 0 // atomic bool: 0 = disabled, 1 = enabled
	listenerCancel context.CancelFunc
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconDisabled)
	systray.SetTitle("D2listener")
	systray.SetTooltip("D2listener")

	startListen := systray.AddMenuItemCheckbox("Listen", "Listen DBUS messages", false)
	mURL := systray.AddMenuItem("Open home page", "my home")
	systray.AddSeparator()
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app")

	go func() {
		for {
			select {
			case <-startListen.ClickedCh:
				if atomic.LoadUint32(&enabled) == 1 {
					// currently enabled -> disable
					atomic.StoreUint32(&enabled, 0)
					systray.SetIcon(iconDisabled)
					startListen.Uncheck()
					if listenerCancel != nil {
						listenerCancel()
					}
					log.Println("Listener disabled")
				} else {
					// disabled -> enable
					atomic.StoreUint32(&enabled, 1)
					systray.SetIcon(iconEnabled)
					startListen.Check()
					ctx, cancel := context.WithCancel(context.Background())
					listenerCancel = cancel
					log.Println("Listener enabled, starting DBus listener")
					go func() {
						if err := listenDBUS(ctx); err != nil {
							log.Printf("DBus listener error: %v\n", err)
						}
					}()
				}
			case <-mURL.ClickedCh:
				if err := open.Run(homePage); err != nil {
					log.Printf("Failed to open home page: %v\n", err)
				}
			case <-mQuit.ClickedCh:
				if listenerCancel != nil {
					listenerCancel()
				}
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleanup happens via context cancellation in goroutines
}

func listenDBUS(ctx context.Context) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}
	defer conn.Close()

	call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"eavesdrop='true',type='method_call',interface='org.freedesktop.Notifications',member='Notify',path='/org/freedesktop/Notifications'")
	if call.Err != nil {
		return fmt.Errorf("failed to add match: %w", call.Err)
	}

	c := make(chan *dbus.Message, 10)
	conn.Eavesdrop(c)
	log.Println("Listening for Dota 2 notifications")

	for {
		select {
		case <-ctx.Done():
			log.Println("DBus listener stopped")
			return nil
		case v := <-c:
			if v == nil {
				continue
			}
			if isGameReadyMessage(v) {
				if atomic.LoadUint32(&enabled) == 1 {
					log.Println("Game ready detected, accepting match")
					if err := findIds(dota2ID); err != nil {
						log.Printf("Error accepting match: %v\n", err)
					}
				}
			}
		}
	}
}

// isGameReadyMessage checks if a DBus message indicates a game ready event
func isGameReadyMessage(msg *dbus.Message) bool {
	return isGameReadyText(msg.String())
}

func isGameReadyText(msg string) bool {
	return strings.Contains(msg, GAMEISREADY) || strings.Contains(msg, GAMEISUNPAUSING)
}
