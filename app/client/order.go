package client

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	d "github.com/denniswon/tcex/app/data"
)

// Client defines typed wrappers for the Ethereum RPC API.
type OrderClient struct {
	file 							*os.File;
	scanner 					*bufio.Scanner;
	mutex 						*sync.RWMutex
	stopped     			bool
	stopChannel 			chan string
	errorChannel		 	chan error
	batchSize 			uint64
	orderNumber 			uint64
	interval    			time.Duration // The interval with which to run the Action
	period      			time.Duration // The actual period of the wait
}

// NewClient creates a client that uses the given RPC client.
func NewOrderClient(interval time.Duration, batchSize uint64) *OrderClient {
	client := &OrderClient{
		stopped:					false,
		stopChannel: 			make(chan string),
		errorChannel: 		make(chan error),
		interval:   			interval,
		period:     			interval,
		batchSize: 				batchSize,
		mutex:      			&sync.RWMutex{},
	}

	initialized := client.Load("trades.txt")
	if !initialized {
		log.Println("Failed to initialize order client")
		return nil
	}

	return client
}

func (c *OrderClient) Load(filename string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.file == nil {
		file, err := os.Open(filename)
		if err != nil {
			log.Println("Error opening file:", err)
			return false
		}
		c.file = file
	} else {
		c.file.Seek(0, 0)	// move to beginning of the file
	}

	c.scanner = bufio.NewScanner(c.file)

	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		return advance, token, err
	}
	c.scanner.Split(scanLines)

	c.orderNumber = 0

	return true
}

func (c *OrderClient) Start(ordersChannel chan d.Orders) {
	log.Println("Worker Started")

	for {
		select {
		case <-c.errorChannel:
			c.Stop()
			break
		case <-c.stopChannel:
			c.stopChannel <- "Stop"
			return
		case <-time.After(c.period):
			// This breaks out of the select, not the for loop.
			break
		}

		started := time.Now()
		c.Receive(ordersChannel)
		finished := time.Now()

		duration := finished.Sub(started)
		c.period = c.interval - duration
	}
}

func (c *OrderClient) Err() chan error {
	return c.errorChannel
}

func (c *OrderClient) Stop() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.stopped = true

	c.stopChannel <- "Stop"
	<-c.stopChannel

	close(c.stopChannel)
	c.file.Close()
}

func (c *OrderClient) IsStopped() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.stopped
}

func (c *OrderClient) Receive(ordersChannel chan d.Orders) {
	orders := []*d.Order{}

	for i:=0; i < int(c.batchSize); i++ {

		order, err := c.Next()

		if err != nil {
			log.Printf("[!] Failed to get next order : %s\n", err.Error())
			c.errorChannel <- err
			continue
		}

		if order == nil {
			log.Println("Done reading trades.txt")
			c.Stop()
			break
		}

		orders = append(orders, order)

	}

	ordersChannel <- d.Orders{Orders: orders}
}

func (c *OrderClient) Next() (*d.Order, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.IsStopped() {
		return nil, nil
	}

	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}

	var order d.Order

	text := c.scanner.Text()
	_text := []byte(text)

	err := json.Unmarshal(_text, &order)
	if err != nil {
		log.Printf("[!] Failed to decode order data to JSON : %s\n", err.Error())
		return nil, err
	}

	order.Number = c.orderNumber
	c.orderNumber++

	return &order, nil
}
