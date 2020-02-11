package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetBidAsc(t *testing.T) {
	input := []byte(`{"e":"depthUpdate","E":1581433169070,"s":"ETHBTC","U":1009388485,"u":1009388515,"b":[["0.02257200","0.71900000"],["0.02255700","90.00000000"]],"a":[["0.02257600","0.23300000"],["0.02258100","0.99100000"]]}`)
	expected := []byte(`{"bid": {"price":0.022572,"amount":0.719}, "ask": {"price": 0.022576, "amount": 0.233}}`)
	res, err := GetBidAsc(input)
	if err != nil && !reflect.DeepEqual(expected, res) {
		t.Error("expected", string(expected), "have", string(res))
	}
	fmt.Println(string(res))
}
