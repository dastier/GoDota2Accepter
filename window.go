//go:build linux

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	maxFindAttempts = 4
	acceptDelay     = 2 * time.Second
)

var (
	xdotoolCmd = "xdotool"
	ydotoolCmd = "ydotool"
)

func findIds(id string) error {
	displayServer := detectDisplayServer()

	for attempt := 1; attempt <= maxFindAttempts; attempt++ {
		log.Printf("Searching for window: %s (attempt %d/%d)", id, attempt, maxFindAttempts)

		wid, err := findWindow(id)
		if err != nil {
			log.Printf("Failed to find window for %s: %v", id, err)
			continue
		}

		log.Printf("Found window ID: %s", wid)
		log.Printf("Window title: %s", getWindowTitle(wid))

		switch displayServer {
		case "x11":
			if err := activateX11(wid); err != nil {
				return fmt.Errorf("could not activate window: %w", err)
			}
		case "wayland":
			log.Println("Wayland: skipping window activation (not supported)")
		default:
			return fmt.Errorf("unsupported display server: %s", displayServer)
		}

		log.Printf("Accepting match in window %s", wid)
		time.Sleep(acceptDelay)
		if err := pressEnter(displayServer, wid); err != nil {
			return fmt.Errorf("could not press enter: %w", err)
		}
		log.Printf("Sent accept keypress to window %s", wid)
		return nil
	}

	return fmt.Errorf("no window found for %s", id)
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

func findWindow(name string) (string, error) {
	out, err := exec.Command(xdotoolCmd, "search", "--name", name).Output()
	if err == nil {
		if wid := strings.TrimSpace(string(out)); wid != "" {
			return strings.SplitN(wid, "\n", 2)[0], nil
		}
	}

	out, err = exec.Command(xdotoolCmd, "search", "--class", name).Output()
	if err == nil {
		if wid := strings.TrimSpace(string(out)); wid != "" {
			return strings.SplitN(wid, "\n", 2)[0], nil
		}
	}

	pids, err := exec.Command("pgrep", "-if", name).Output()
	if err != nil {
		return "", fmt.Errorf("no process found for %s", name)
	}
	for _, pid := range strings.Fields(string(pids)) {
		out, err := exec.Command(xdotoolCmd, "search", "--pid", pid).Output()
		if err == nil {
			if wid := strings.TrimSpace(string(out)); wid != "" {
				return strings.SplitN(wid, "\n", 2)[0], nil
			}
		}
	}

	return "", fmt.Errorf("no window found for %s", name)
}

func getWindowTitle(wid string) string {
	out, err := exec.Command(xdotoolCmd, "getwindowname", wid).Output()
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	return strings.TrimSpace(string(out))
}

func activateX11(wid string) error {
	if out, err := exec.Command(xdotoolCmd, "windowactivate", "--sync", wid).CombinedOutput(); err != nil {
		return fmt.Errorf("windowactivate failed: %w, output: %s", err, strings.TrimSpace(string(out)))
	}

	log.Printf("Activated window %s", wid)
	if out, err := exec.Command(xdotoolCmd, "windowsize", wid, "100%", "100%").CombinedOutput(); err != nil {
		log.Printf("Warning: windowsize failed: %v, output: %s", err, strings.TrimSpace(string(out)))
	} else {
		log.Printf("Maximized window %s", wid)
	}
	return nil
}

func pressEnter(displayServer, wid string) error {
	switch displayServer {
	case "x11":
		return exec.Command(xdotoolCmd, "key", "--window", wid, "Return").Run()
	case "wayland":
		return exec.Command(ydotoolCmd, "key", "28:1", "28:0").Run()
	default:
		return fmt.Errorf("unsupported display server: %s", displayServer)
	}
}
