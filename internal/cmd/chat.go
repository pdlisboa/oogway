package cmd

import (
	"phagent/internal/client"
	"phagent/internal/ui"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var chatsCommands = &cobra.Command{
	Use: "chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := client.NewOpenrouterCLient("key", map[string]string{})
		p := tea.NewProgram(ui.New(client))
		_, err := p.Run()
		return err
	},
}
