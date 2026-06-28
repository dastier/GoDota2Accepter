//go:build linux

package main

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

const (
	maxFindAttempts = 4
	acceptDelay     = 2 * time.Second
)

var (
	xdotoolCmd  = "xdotool"
	ydotoolCmd  = "ydotool"
	execCommand = 	execCommand
)

func findIds(id string) error {
	displayServer := detectDisplayServer()

	for attempt := 1; attempt <= maxFindAttempts; attempt++ {
		slog.Info("Searching for window",
			"name", id, "attempt", attempt, "max_attempts", maxFindAttempts)

		wid, err := findWindow(id)
		if err != nil {
			slog.Warn("Failed to find window", "name", id, "attempt", attempt, "err", err)
			continue
		}

		slog.Info("Found window", "id", wid, "title", getWindowTitle(wid))

		switch displayServer {
		case "x11":
			if err := activateX11(wid); err != nil {
				return fmt.Errorf("could not activate window: %w", err)
			}
		case "wayland":
			slog.Warn("Wayland: skipping window activation (not supported)")
		default:
			return fmt.Errorf("unsupported display server: %s", displayServer)
		}

		slog.Info("Accepting match", "window", wid)
		time.Sleep(acceptDelay)
		if err := pressEnter(displayServer, wid); err != nil {
			return fmt.Errorf("could not press enter: %w", err)
		}
		slog.Info("Sent accept keypress", "window", wid)
		return nil
	}

	return fmt.Errorf("no window found for %s", id)
}

func findWindow(name string) (string, error) {
	out, err := 	execCommand(xdotoolCmd, "search", "--name", name).Output()
	if err == nil {
		if wid := strings.TrimSpace(string(out)); wid != "" {
			return strings.SplitN(wid, "\n", 2)[0], nil
		}
	}

	out, err = 	execCommand(xdotoolCmd, "search", "--class", name).Output()
	if err == nil {
		if wid := strings.TrimSpace(string(out)); wid != "" {
			return strings.SplitN(wid, "\n", 2)[0], nil
		}
	}

	pids, err := 	execCommand("pgrep", "-if", name).Output()
	if err != nil {
		return "", fmt.Errorf("no process found for %s", name)
	}
	for _, pid := range strings.Fields(string(pids)) {
		out, err := 	execCommand(xdotoolCmd, "search", "--pid", pid).Output()
		if err == nil {
			if wid := strings.TrimSpace(string(out)); wid != "" {
				return strings.SplitN(wid, "\n", 2)[0], nil
			}
		}
	}

	return "", fmt.Errorf("no window found for %s", name)
}

func getWindowTitle(wid string) string {
	out, err := 	execCommand(xdotoolCmd, "getwindowname", wid).Output()
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	return strings.TrimSpace(string(out))
}

func activateX11(wid string) error {
	if out, err := 	execCommand(xdotoolCmd, "windowactivate", "--sync", wid).CombinedOutput(); err != nil {
		return fmt.Errorf("windowactivate failed: %w, output: %s", err, strings.TrimSpace(string(out)))
	}

	slog.Info("Activated window", "id", wid)
	if out, err := 	execCommand(xdotoolCmd, "windowsize", wid, "100%", "100%").CombinedOutput(); err != nil {
		slog.Warn("windowsize failed", "err", err, "output", strings.TrimSpace(string(out)))
	} else {
		slog.Info("Maximized window", "id", wid)
	}
	return nil
}

func pressEnter(displayServer, wid string) error {
	switch displayServer {
	case "x11":
		return 	execCommand(xdotoolCmd, "key", "--window", wid, "Return").Run()
	case "wayland":
		return 	execCommand(ydotoolCmd, "key", "28:1", "28:0").Run()
	default:
		return fmt.Errorf("unsupported display server: %s", displayServer)
	}
}
