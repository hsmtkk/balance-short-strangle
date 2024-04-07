package put

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/boardget"
	"github.com/hsmtkk/aukabucomgo/info/symbolnameoptionget"
	"github.com/hsmtkk/balance-short-strangle/config"
	"github.com/hsmtkk/balance-short-strangle/future"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use: "put",
	Run: run,
}

var targetDelta float64
var year int
var month int

func init() {
	nextMonth := time.Now().AddDate(0, 1, 0)
	Command.Flags().Float64Var(&targetDelta, "delta", config.TARGET_DELTA, "target delta")
	Command.Flags().IntVar(&year, "year", nextMonth.Year(), "year")
	Command.Flags().IntVar(&month, "month", int(nextMonth.Month()), "month")
}

func run(command *cobra.Command, args []string) {
	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		log.Fatal("env var API_PASSWORD is not defined")
	}
	baseClient, err := base.New(base.PRODUCTION, apiPassword)
	if err != nil {
		log.Fatal(err)
	}
	futurePrice, err := future.Price(baseClient, year, month)
	if err != nil {
		log.Fatal(err)
	}
	symbolName, delta, err := findOption(baseClient, futurePrice, year, month, targetDelta)
	if err != nil {
		log.Fatal(err)
	}
	displayResult(symbolName, delta)
}

func findOption(baseClient base.Client, futurePrice, year, month int, targetDelta float64) (string, float64, error) {
	boardClient := boardget.New(baseClient)
	symbolClient := symbolnameoptionget.New(baseClient)
	strikePrice := (futurePrice / 250) * 250
	for {
		symbolResp, err := symbolClient.SymbolNameOptionGet(symbolnameoptionget.NK225miniop, year, month, symbolnameoptionget.PUT, strikePrice)
		if err != nil {
			return "", 0, err
		}
		symbol := symbolResp.Symbol
		boardResp, err := boardClient.BoardGet(symbol, boardget.ALL_DAY)
		if err != nil {
			return "", 0, err
		}
		if boardResp.Delta > targetDelta {
			return boardResp.SymbolName, boardResp.Delta, nil
		}
		strikePrice -= 250
	}
	return "", 0, fmt.Errorf("failed to find matching option")
}

func displayResult(symbolName string, delta float64) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"symbol name", "delta"})
	t.AppendRow([]interface{}{symbolName, delta})
	t.Render()
}
