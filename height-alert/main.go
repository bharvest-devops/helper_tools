package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cometbftHttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cometbft/cometbft/types"
)

func timeCalc(sec uint64) (uint64, uint64, uint64, uint64) {
	var remainder uint64 = sec
	seconds := remainder % 60
	remainder /= 60

	minutes := remainder % 60
	remainder /= 60

	hours := remainder % 24
	days := uint64(remainder / 24)

	return days, hours, minutes, seconds
}

func sendMsg(token, channel, msg string) {
	req := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, channel, msg)
	_, err := http.Get(req)
	if err != nil {
		panic(err)
	}
}

func main() {
	chain := "Cosmos Mainnet"
	var upgradeHeight int64 = 591800
	blockTime := 5.0
	rpcHost := "http://localhost:26657"

	tgToken := ""
	tgChannel := ""

	client, err := cometbftHttp.New(rpcHost, "/websocket")
	if err != nil {
		panic(err)
	}
	client.Start()
	blockHeader, err := client.WSEvents.Subscribe(context.Background(), "", "tm.event = 'NewBlockHeader'")
	defer client.Stop()
	if err != nil {
		panic(err)
	}
	fmt.Println("Start")
	for e := range blockHeader {
		var currentHeight int64 = e.Data.(types.EventDataNewBlockHeader).Header.Height
		var remainSeconds uint64 = uint64(blockTime * float64(upgradeHeight-currentHeight))
		days, hours, minutes, seconds := timeCalc(remainSeconds)

		if (days == 1 && hours == 0 && minutes == 0) || (days == 0 && hours == 1 && minutes == 0) || (days == 0 && hours == 0 && minutes == 0 && seconds == 0) {
			msg := fmt.Sprintf("Remain time: %dDays %dHours %dMinutes %dSeconds", days, hours, minutes, seconds)
			sendMsg(tgToken, tgChannel, "🔴 "+chain+" Upgrade 🔴%0A"+msg)
			time.Sleep(70 * time.Second)
		}

		fmt.Printf("\nCurrent height / Upgrade height : %d / %d\n", currentHeight, upgradeHeight)
		fmt.Printf("Remain time: %dDays %dHours %dMinutes %dSeconds\n", days, hours, minutes, seconds)

	}
}
