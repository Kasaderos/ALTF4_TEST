package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	names = []string{"btcusdt", "ethusdt", "eosusdt"}
)

type Output struct {
	Bid Pair `json:"bid"`
	Ask Pair `json:"ask"`
}

type Pair struct {
	Price  float64 `json:"price","float64"`
	Amount float64 `json:"amount","float64"`
}

func ParseJSON(key string, in interface{}) ([]Pair, error) {
	res := in.(map[string]interface{})
	if _, ok := res[key]; !ok {
		return nil, fmt.Errorf("ParseJSON: key error")
	}
	listOfBids := res[key].([]interface{})
	var listOfPair []Pair
	for _, p := range listOfBids {
		pair := p.([]interface{})
		price, err := strconv.ParseFloat(pair[0].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("strconv: price")
		}
		amount, err := strconv.ParseFloat(pair[1].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("strconv: amount")
		}
		p := Pair{Price: price, Amount: amount}
		listOfPair = append(listOfPair, p)
	}
	return listOfPair, nil
}

func FindMax(pairs []Pair) Pair {
	if len(pairs) < 1 {
		return Pair{}
	}
	max := 0
	for i := 1; i < len(pairs); i++ {
		if pairs[i].Price > pairs[max].Price {
			max = i
		}
	}
	return pairs[max]
}

func FindMin(pairs []Pair) Pair {
	if len(pairs) < 1 {
		return Pair{}
	}
	min := 0
	for i := 1; i < len(pairs); i++ {
		if pairs[i].Price < pairs[min].Price {
			min = i
		}
	}
	return pairs[min]
}

func GetBidAsc(data []byte) ([]byte, error) {
	var temp interface{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return nil, fmt.Errorf("json: can't unmarshal")
	}
	pjBid, err := ParseJSON("b", temp)
	if err != nil {
		return nil, fmt.Errorf("ParseJSON: error json string")
	}
	pjAsk, err := ParseJSON("a", temp)
	if err != nil {
		return nil, fmt.Errorf("ParseJSON: error json string")
	}
	pb := FindMax(pjBid)
	pa := FindMin(pjAsk)
	out := Output{Bid: pb, Ask: pa}
	result, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("json: can't marshal res")
	}
	return result, nil
}

func GetFromBinance(name string, timeout time.Duration) {
	u := url.URL{
		Scheme: "wss",
		Host:   "stream.binance.com:9443",
		Path:   fmt.Sprintf("/ws/%s@depth", name),
	}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c.SetReadDeadline(time.Now().Add(timeout))
	defer func() {
		err = c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// можно вынести в горутину
		out, err := GetBidAsc(message)

		if err != nil {
			log.Print(err)
			return
		}
		log.Println(fmt.Sprintf("%s %s", name, string(out)))
		//log.Println(string(message))
	}

}

func main() {
	wg := &sync.WaitGroup{}
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			GetFromBinance(name, 3*time.Second)
		}(name)
	}
	wg.Wait()
}
