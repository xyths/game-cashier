package puller

import "testing"

func TestToAmount(t *testing.T) {
	type testpair struct {
		Quantity string
		Amount   float64
	}
	testpairs := []testpair{
		{"0.0020 EOS", 0.002},
		{"", 0},
		{"xx", 0},
		{"1000.00 ABC", 0},
		{"100.100 EOS", 100.1},
	}
	for _, pair := range testpairs {
		got := toAmount(pair.Quantity)
		if got != pair.Amount {
			t.Errorf("quantity: %s, expect: %f, got: %f", pair.Quantity, pair.Amount, got)
		}
	}
}
