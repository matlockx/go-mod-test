package go-mod-test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reStringAmount = regexp.MustCompile(`^([0-9]+)(?:[.,]([0-9]{0,2}))? ([A-Z]{3})$`)
var reCurrency = regexp.MustCompile("^[A-Z]{3}$")
var reAmountInMajor = regexp.MustCompile("^[0-9]+[.,][0-9]{2}$")

func isCurrencyCode(code string) bool {
	return reCurrency.MatchString(code)
}

func (m Amount) MarshalJSON() ([]byte, error) {
	var str string

	if m.currency.isDefined() {
		str = fmt.Sprintf(`{"amountInMinor":%d,"amountInMajor":%1.2f,"amount":"%1.2f","currency":"%s"}`,
			m.amountInMinor, m.InMajorUnits(), m.InMajorUnits(), m.currency.code)
	} else {
		str = "null"
	}

	return []byte(str), nil
}

func (m *Amount) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		m.amountInMinor = 0
		m.currency = Currency{}
		return nil
	}

	var decoded interface{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return err
	}

	if decoded, ok := decoded.(string); ok {
		match := reStringAmount.FindStringSubmatch(decoded)
		if match == nil {
			return fmt.Errorf("invalid money amount: '%s'", decoded)
		}

		amountInMinor, err := strconv.Atoi(match[1] + match[2] + ("00"[len(match[2]):]))
		if err != nil {
			return fmt.Errorf("could not parse amount in minor")
		}

		m.amountInMinor = amountInMinor
		m.currency = Currency{match[3]}
		return nil
	}

	if decoded, ok := decoded.(map[string]interface{}); ok {
		currency, ok := decoded["currency"].(string)
		if !ok || !isCurrencyCode(currency) {
			return fmt.Errorf("invalid currency given: '%s'", currency)
		}

		// check if we have a amountInMinor property, use this with priority.
		if amountInMinor, ok := decoded["amountInMinor"].(json.Number); ok {
			amountInMinor, _ := amountInMinor.Int64()
			m.amountInMinor = int(amountInMinor)
			m.currency = Currency{currency}
			return nil
		}

		amount := decoded["amount"]
		if amount == nil {
			m.amountInMinor = 0
			m.currency = Currency{currency}
			return nil
		}

		if amount, ok := amount.(json.Number); ok {
			amountInMinor, err := strconv.Atoi(string(amount))
			if err != nil {
				return fmt.Errorf("could not parse amount in minor")
			}

			m.amountInMinor = amountInMinor
			m.currency = Currency{currency}
			return nil
		}

		if amount, ok := amount.(string); ok {
			if reAmountInMajor.MatchString(amount) {
				amountInMinorString := strings.Replace(strings.Replace(amount, ".", "", 1), ",", "", 1)

				amountInMinor, err := strconv.Atoi(amountInMinorString)
				if err != nil {
					return fmt.Errorf("could not parse amount in minor")
				}

				m.amountInMinor = amountInMinor
				m.currency = Currency{currency}
				return nil
			}

			// parse it directly in minor
			amountInMinor, err := strconv.Atoi(amount)
			if err != nil {
				return fmt.Errorf("could not parse amount in minor")
			}

			m.amountInMinor = amountInMinor
			m.currency = Currency{currency}
			return nil
		}
	}

	return errors.New("invalid format for money.Amount")
}
