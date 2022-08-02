package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/liamg/guerrilla/internal/app/output"

	"github.com/liamg/guerrilla/pkg/guerrilla"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "guerrilla",
	Short: "Guerrilla is a command line tool for creating a temporary email address and receiving associated email in the terminal. Powered by the Guerrilla Mail API.",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		client, err := guerrilla.Init()
		if err != nil {
			return err
		}

		if flagPollIntervalSeconds < 1 || flagPollIntervalSeconds > 600 {
			return fmt.Errorf("poll-interval must be between 1-600")
		}

		printer := output.New(cmd.OutOrStdout())

		printer.PrintSummary(client.GetAddress())

		poller := guerrilla.NewPoller(client, guerrilla.PollOptionWithInterval(time.Second*time.Duration(flagPollIntervalSeconds)))
		var count int
		for email := range poller.Poll() {
			if !showWelcomeEmail && count == 0 && email.Subject == "Welcome to Guerrilla Mail" {
				continue
			}
			printer.PrintEmail(email)
			count++
		}

		return nil
	},
}

var (
	flagPollIntervalSeconds int
	showWelcomeEmail        bool
)

func Execute() {

	rootCmd.Flags().IntVarP(&flagPollIntervalSeconds, "poll-interval", "i", 30, "Poll interval in seconds. Must be between 1-600. Low values are not recommended due to API rate limits.")
	rootCmd.Flags().BoolVarP(&showWelcomeEmail, "show-welcome", "w", false, "Show the default GuerrillaMail welcome email in the output (filtered by default).")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
