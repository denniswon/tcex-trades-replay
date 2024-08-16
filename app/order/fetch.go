package order

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	d "github.com/denniswon/tcex/app/data"
	q "github.com/denniswon/tcex/app/queue"

	"gorm.io/gorm"
)

// FetchOrderByNumber - Fetching the next order to process
func FetchOrderByNumber(number uint64, _db *gorm.DB, redis *d.RedisInfo, queue *q.OrderReplayQueue, _status *d.StatusHolder) bool {

	// Starting order processing at
	startingAt := time.Now().UTC()

	_num := big.NewInt(0)
	_num.SetUint64(number)


	input, err := os.Open("trades.txt")
	if err != nil {
		log.Println("Error opening file:", err)
		return false
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	i := big.NewInt(0)
	var order: Order{}
	for {
		if i.Cmp(_num) == 0 {
			break
		}

		if !scanner.Scan() || (err := scanner.Err() && err != nil) {
				log.Println("Error reading from trades.txt:", err)
				return false
			}
		}

		line = scanner.Text()
		err := json.Unmarshal([]byte(line), data)
		if err != nil {
				log.Println("Error unmarshalling order", err.Error())
		}
	}

	return ProcessOrderContent(order, _db, redis, true, queue, _status, startingAt)
}
