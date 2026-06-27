//go:build !linux

package main

import "fmt"

func findIds(id string) error {
	return fmt.Errorf("window activation is only supported on Linux: %s", id)
}
