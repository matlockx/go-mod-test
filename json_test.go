package gomodtest

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		in  Amount
		out string
	}{
		{FromMinor(50, Currency{"BTC"}), `{"amountInMinor":50,"amountInMajor":0.5,"amount":"0.50","currency":"BTC"}`},
		{FromMinor(100, Currency{"EUR"}), `{"amountInMinor":100,"amountInMajor":1.0,"amount":"1.00","currency":"EUR"}`},
		{FromMinor(250, Currency{"USD"}), `{"amountInMinor":250,"amountInMajor":2.5,"amount":"2.50","currency":"USD"}`},

		{FromMinor(0, Currency{"USD"}), `{"amountInMinor":0,"amountInMajor":0,"amount":"0.00","currency":"USD"}`},
		{Amount{}, `null`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.in.String(), func(t *testing.T) {
			result, _ := json.Marshal(testCase.in)

			var expected, deserialized map[string]interface{}
			if err := json.Unmarshal([]byte(testCase.out), &expected); err != nil {
				t.Fatal(err)
			}

			if err := json.Unmarshal(result, &deserialized); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(expected, deserialized) {
				t.Errorf("Expected %s, got %s", expected, deserialized)
			}
		})
	}
}

func TestMoneyUnmarshalFromFloat(t *testing.T) {
	testCases := []struct {
		out Amount
		in  string
	}{
		{FromMinor(50, Currency{"BTC"}), `"0.50 BTC"`},
		{FromMinor(150, Currency{"GBP"}), `"1.50 GBP"`},
		{FromMinor(0, Currency{"GBP"}), `"0 GBP"`},
		{FromMinor(20, Currency{"GBP"}), `"0.2 GBP"`},
		{FromMinor(14450, Currency{"CHF"}), `"144.5 CHF"`},

		{FromMinor(50, Currency{"EUR"}), `{"amount":"50","currency":"EUR"}`},
		{FromMinor(100, Currency{"EUR"}), `{"amount":"100","currency":"EUR"}`},
		{FromMinor(250, Currency{"USD"}), `{"amount":"250","currency":"USD"}`},

		{FromMinor(50, Currency{"EUR"}), `{"amount":50,"currency":"EUR"}`},
		{FromMinor(100, Currency{"EUR"}), `{"amount":100,"currency":"EUR"}`},
		{FromMinor(250, Currency{"USD"}), `{"amount":250,"currency":"USD"}`},

		{FromMinor(50, Currency{"EUR"}), `{"amount":"0.50","currency":"EUR"}`},
		{FromMinor(200, Currency{"EUR"}), `{"amount":"2.00","currency":"EUR"}`},
		{FromMinor(185, Currency{"USD"}), `{"amount":"1.85","currency":"USD"}`},

		{FromMinor(50, Currency{"EUR"}), `{"amountInMinor":50,"currency":"EUR"}`},
		{FromMinor(100, Currency{"EUR"}), `{"amountInMinor":100,"currency":"EUR"}`},
		{FromMinor(250, Currency{"USD"}), `{"amountInMinor":250,"currency":"USD"}`},

		{FromMinor(0, Currency{"CHF"}), `{"amount":"0","currency":"CHF"}`},

		{FromMinor(0, Currency{"BTC"}), `{"currency":"BTC"}`},

		{Amount{}, `null`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.in, func(t *testing.T) {
			var result Amount

			err := json.Unmarshal([]byte(testCase.in), &result)
			if err != nil {
				t.Fatalf("Got error %s", err)
			}

			if result != testCase.out {
				t.Errorf("Expected %s, got %s", testCase.out, result)
			}
		})
	}
}
