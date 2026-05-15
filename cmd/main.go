package main

import (
	"asterminal/internal/config"
	"asterminal/internal/model"
	"asterminal/internal/render"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := config.GetConfig()
	model := model.NewModel(config)
	_, err := tea.NewProgram(render.NewRenderer(model)).Run()

	if err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}
