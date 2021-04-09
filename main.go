package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus/v5"
	"github.com/skratchdot/open-golang/open"
)

const (
	IconEnabled     = "assets/d_enabled.svg"
	IconDisabled    = "assets/d_disabled.svg"
	GAMEISREADY     = "Your game is ready"
	GAMEISUNPAUSING = "The game is unpausing..."
	dota2ID         = "Dota2"
	homePage        = "https://github.com/dastier/GoDota2Accepter"
)

var ENABLED bool

func main() {
	ENABLED = false
	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetIcon(getIcon(IconDisabled))
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
				if startListen.Checked() {
					ENABLED = false
					systray.SetIcon(getIcon(IconDisabled))
					startListen.Uncheck()
					fmt.Println("false and uncheck")
				} else {
					ENABLED = true
					systray.SetIcon(getIcon(IconEnabled))
					startListen.Check()
					fmt.Println("true and check and launch")
					go listenDBUS()
				}
			case <-mURL.ClickedCh:
				err := open.Run(homePage)
				if err != nil {
					os.Exit(1)
				}

			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
}

func listenDBUS() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"eavesdrop='true',type='method_call',interface='org.freedesktop.Notifications',member='Notify',path='/org/freedesktop/Notifications'")
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "Failed to add match:", call.Err)
		os.Exit(1)
	}
	c := make(chan *dbus.Message, 10)
	conn.Eavesdrop(c)
	fmt.Println("Listening for Dota messages")

	for v := range c {
		if (strings.Contains(v.String(), GAMEISREADY)) || (strings.Contains(v.String(), GAMEISUNPAUSING)) {
			if ENABLED {
				fmt.Println("FOUND!")
				fmt.Println("call finIds from loop")
				findIds(dota2ID)
			}
		}
	}
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
