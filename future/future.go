package future

import (
	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/boardget"
	"github.com/hsmtkk/aukabucomgo/info/symbolnamefutureget"
)

func Price(baseClient base.Client, year, month int) (int, error) {
	symbolClient := symbolnamefutureget.New(baseClient)
	symbolResp, err := symbolClient.SymbolNameFutureGet(symbolnamefutureget.NK225mini, year, month)
	if err != nil {
		return 0, err
	}
	symbol := symbolResp.Symbol

	boardClient := boardget.New(baseClient)
	boardResp, err := boardClient.BoardGet(symbol, boardget.ALL_DAY)
	if err != nil {
		return 0, err
	}
	return int(boardResp.CurrentPrice), nil
}
