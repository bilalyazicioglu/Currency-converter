package cmd

import (
	"fmt"
	"os"

	"Currency-Converter/internal/api"
	"Currency-Converter/internal/config"
	"Currency-Converter/internal/ui"

	"github.com/spf13/cobra"
)

var (
	fromCurrency string
	toCurrency   string
	amount       float64
)

var rootCmd = &cobra.Command{
	Use:   "currency-converter",
	Short: "A terminal-based currency converter",
	Long: `Currency Converter is a beautiful terminal application
that converts currencies using real-time exchange rates.

Built with Bubble Tea (TUI) and Cobra CLI.
Powered by ExchangeRate-API.

Examples:
  currency-converter              # Start interactive TUI
  currency-converter -f USD -t EUR -a 100   # Quick conversion`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client := api.NewClient(cfg.APIKey)

		if fromCurrency != "" && toCurrency != "" && amount > 0 {
			result, err := client.Convert(fromCurrency, toCurrency, amount)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\nCurrency Conversion:\n\n")
			fmt.Printf("  %.2f %s = %.2f %s\n", result.Amount, result.FromCurrency, result.Result, result.ToCurrency)
			fmt.Printf("  Rate: 1 %s = %.4f %s\n", result.FromCurrency, result.Rate, result.ToCurrency)
			return
		}

		if err := ui.Run(client); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&fromCurrency, "from", "f", "", "Source currency code (e.g., USD)")
	rootCmd.Flags().StringVarP(&toCurrency, "to", "t", "", "Target currency code (e.g., EUR)")
	rootCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Amount to convert")
}
