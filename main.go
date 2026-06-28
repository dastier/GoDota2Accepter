package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"

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
	dbusChanBufSize = 10
)

type listener struct {
	cancel context.CancelFunc
}

var state atomic.Pointer[listener]

//go:embed assets/d_enabled.svg
var iconEnabled []byte

//go:embed assets/d_disabled.svg
var iconDisabled []byte

func main() {
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh
		if old := state.Swap(nil); old != nil {
			old.cancel()
		}
		systray.Quit()
	}()
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
				if old := state.Swap(nil); old != nil {
					old.cancel()
					systray.SetIcon(iconDisabled)
					startListen.Uncheck()
					slog.Info("Listener disabled")
				} else {
					systray.SetIcon(iconEnabled)
					startListen.Check()
					ctx, cancel := context.WithCancel(context.Background())
					state.Store(&listener{cancel: cancel})
					slog.Info("Listener enabled, starting DBus listener")
					go func() {
						if err := listenDBUS(ctx); err != nil {
							slog.Error("DBus listener error", "err", err)
						}
					}()
				}
			case <-mURL.ClickedCh:
				if err := open.Run(homePage); err != nil {
					slog.Error("Failed to open home page", "err", err)
				}
			case <-mQuit.ClickedCh:
				if old := state.Swap(nil); old != nil {
					old.cancel()
				}
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {}

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

	c := make(chan *dbus.Message, dbusChanBufSize)
	conn.Eavesdrop(c)
	slog.Info("Listening for Dota 2 notifications")

	for {
		select {
		case <-ctx.Done():
			slog.Info("DBus listener stopped")
			return nil
		case v := <-c:
			if v == nil {
				continue
			}
			if isGameReadyText(v.String()) && state.Load() != nil {
				slog.Info("Game ready detected, accepting match")
				if err := findIds(dota2ID); err != nil {
					slog.Error("Error accepting match", "err", err)
				}
			}
		}
	}
}

func detectDisplayServer() string {
	if v := os.Getenv("WAYLAND_DISPLAY"); v != "" {
		return "wayland"
	}
	if v := os.Getenv("XDG_SESSION_TYPE"); v == "wayland" {
		return "wayland"
	}
	return "x11"
}

func isGameReadyText(msg string) bool {
	return strings.Contains(msg, GAMEISREADY) || strings.Contains(msg, GAMEISUNPAUSING)
}
