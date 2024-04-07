package command

import (
	"fmt"
	"log"
	"os"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/boardget"
	"github.com/hsmtkk/aukabucomgo/info/positionsget"
	"github.com/hsmtkk/balance-short-strangle/command/call"
	"github.com/hsmtkk/balance-short-strangle/command/put"
	"github.com/hsmtkk/balance-short-strangle/config"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use: "balance-short-strangle",
	Run: run,
}

var targetDelta float64

func init() {
	Command.AddCommand(call.Command)
	Command.AddCommand(put.Command)
	Command.Flags().Float64Var(&targetDelta, "delta", config.TARGET_DELTA, "target delta")
}

func run(command *cobra.Command, args []string) {
	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		log.Fatal("env var API_PASSWORD is not defined")
	}
	pos, err := getCurrentShortPosition(apiPassword)
	if err != nil {
		log.Fatal(err)
	}
	displayCurrentShortPosition(pos)
	displayTotalDelta(pos)
}

type shortPosition []shortOption

type shortOption struct {
	symbolName string
	quantity   int
	delta      float64
}

func getCurrentShortPosition(apiPassword string) (shortPosition, error) {
	baseClient, err := base.New(base.PRODUCTION, apiPassword)
	if err != nil {
		return nil, err
	}
	positionClient := positionsget.New(baseClient)
	positions, err := positionClient.PositionsGet(positionsget.OPTION, positionsget.SELL)
	if err != nil {
		return nil, err
	}
	boardClient := boardget.New(baseClient)
	result := []shortOption{}
	for _, pos := range positions {
		board, err := boardClient.BoardGet(pos.Symbol, boardget.ALL_DAY)
		if err != nil {
			return nil, err
		}
		result = append(result, shortOption{
			symbolName: pos.SymbolName,
			quantity:   int(pos.LeavesQty),
			delta:      board.Delta,
		})
	}
	return result, nil
}

func displayCurrentShortPosition(position shortPosition) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"symbol name", "quantity", "delta"})
	for _, opt := range position {
		t.AppendRow([]interface{}{opt.symbolName, opt.quantity, opt.delta})
	}
	t.Render()
}

func displayTotalDelta(position shortPosition) {
	totalDelta := 0.0
	for _, opt := range position {
		totalDelta += -float64(opt.quantity) * opt.delta // 売りなのでマイナス
	}
	fmt.Printf("Total delta: %f\n", totalDelta)
	if totalDelta > targetDelta {
		fmt.Println("Sell CALL for delta neutral")
	} else if totalDelta < -targetDelta {
		fmt.Println("Sell PUT for delta neutral")
	} else {
		fmt.Println("Delta neutral")
	}
}
