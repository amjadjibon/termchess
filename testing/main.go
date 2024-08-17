package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func main() {
	var country string
	s := huh.NewSelect[string]().
		Title("choose a piece").
		Options(
			huh.NewOption("Queen", "Q"),
			huh.NewOption("Rook", "R"),
			huh.NewOption("Bishop", "B"),
			huh.NewOption("Knight", "N"),
		).
		Value(&country)

	err := s.Run()
	if err != nil {
		panic(err)
	}

	var v string
	s.Value(&v)

	fmt.Println(v)
}
