//go:build linux

package main

import (
	"os/exec"
	"testing"
)

func TestFindWindow(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "1234")
	}
	defer func() { execCommand = exec.Command }()

	wid, err := findWindow("Dota2")
	if err != nil {
		t.Fatalf("findWindow() = _, %v", err)
	}
	if wid != "1234" {
		t.Fatalf("findWindow() = %q, want %q", wid, "1234")
	}
}

func TestFindWindowNotFound(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { execCommand = exec.Command }()

	_, err := findWindow("NonExistent")
	if err == nil {
		t.Fatal("findWindow() expected error, got nil")
	}
}

func TestActivateX11(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { execCommand = exec.Command }()

	if err := activateX11("1234"); err != nil {
		t.Fatalf("activateX11() = %v", err)
	}
}

func TestPressEnterX11(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { execCommand = exec.Command }()

	if err := pressEnter("x11", "1234"); err != nil {
		t.Fatalf("pressEnter() = %v", err)
	}
}

func TestPressEnterWayland(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { execCommand = exec.Command }()

	if err := pressEnter("wayland", "1234"); err != nil {
		t.Fatalf("pressEnter() = %v", err)
	}
}

func TestPressEnterUnsupported(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { execCommand = exec.Command }()

	if err := pressEnter("unsupported", "1234"); err == nil {
		t.Fatal("pressEnter() expected error, got nil")
	}
}
