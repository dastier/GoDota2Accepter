//go:build linux

package main

import (
	"fmt"
	"log"

	"github.com/go-vgo/robotgo"
)

const maxFindAttempts = 4

func findIds(id string) error {
	for attempt := 1; attempt <= maxFindAttempts; attempt++ {
		log.Printf("Searching for process IDs by name: %s (attempt %d/%d)", id, attempt, maxFindAttempts)
		fpid, err := robotgo.FindIds(id)
		if err != nil {
			log.Printf("Failed to find process IDs for %s: %v", id, err)
			return fmt.Errorf("found errors while robotgo.FindIds: %w", err)
		}

		log.Printf("Found process IDs for %s: %v", id, fpid)
		if len(fpid) == 0 {
			continue
		}

		pid := fpid[0]
		if err := robotgo.ActivePID(pid); err != nil {
			return fmt.Errorf("could not activate dota window: %w", err)
		}

		title := robotgo.GetTitle(pid)
		log.Printf("Activated process %d with window title %q", pid, title)

		x, y, w, h := robotgo.GetBounds(pid)
		log.Printf("Window bounds for process %d: x=%d y=%d width=%d height=%d", pid, x, y, w, h)
		robotgo.MaxWindow(pid)
		robotgo.Sleep(2)
		robotgo.KeyTap("enter")
		log.Printf("Sent accept keypress to process %d", pid)
		return nil
	}

	return fmt.Errorf("no process IDs found for %s", id)
}
