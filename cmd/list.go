package cmd

import (
	"fmt"
	"os"
	"remotelink/config"
	"remotelink/models"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var listFancyCmd = &cobra.Command{
	Use:   "ls",
	Short: "List servers",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialFancyModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	cellStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(20)

	selectedCellStyle = cellStyle.Copy().
				Background(lipgloss.Color("#3C3C3C")).
				Bold(true)

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

type fancyModel struct {
	cursor  int
	servers []models.Server
}

func initialFancyModel() fancyModel {
	return fancyModel{
		cursor:  0,
		servers: config.Servers,
	}
}

func (m fancyModel) Init() tea.Cmd {
	return nil
}

func (m fancyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.servers)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m fancyModel) View() string {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render("Name") + " ")
	s.WriteString(headerStyle.Render("Host") + " ")
	s.WriteString(headerStyle.Render("User") + " ")
	s.WriteString(headerStyle.Render("Containers"))
	s.WriteString("\n")

	s.WriteString(borderStyle.Render(strings.Repeat("─", 84)))
	s.WriteString("\n")

	// Rows
	for i, server := range m.servers {
		style := cellStyle
		if i == m.cursor {
			style = selectedCellStyle
		}

		cursor := "  "
		if i == m.cursor {
			cursor = "▶ "
		}

		s.WriteString(cursor)
		s.WriteString(style.Render(truncate(server.ServerName, 18)))
		s.WriteString(" ")
		s.WriteString(style.Render(fmt.Sprintf("%s:%d", server.HostIp, server.Port)))
		s.WriteString(" ")
		s.WriteString(style.Render(truncate(server.Username, 18)))
		s.WriteString(" ")
		s.WriteString(style.Render(fmt.Sprintf("%d", len(server.Containers))))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/↓: navigate • q: quit"))

	return s.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func init() {
	rootCmd.AddCommand(listFancyCmd)
}
