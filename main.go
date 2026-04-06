package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/mrjira/export"
)

func main() {
	m := export.NewApp()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
