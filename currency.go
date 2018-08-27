package go-mod-test

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

var ErrInvalidCurrencyCode = errors.New("illegal currency code")

type Currency struct {
	code string
}

func (c Currency) String() string {
	return c.code
}

// Parses the input into a three letter currency code. Fails,
// if the input can not be parsed to a currency code.
func ParseCurrency(input string) (Currency, error) {
	if !isCurrencyCode(input) {
		return Currency{}, ErrInvalidCurrencyCode
	}

	return Currency{input}, nil
}

func MustParseCurrency(input string) Currency {
	currency, err := ParseCurrency(input)
	if err != nil {
		panic(err)
	}

	return currency
}

func (c *Currency) isDefined() bool {
	return len(c.code) != 0
}

func (c Currency) Value() (driver.Value, error) {
	return c.code, nil
}

func (c *Currency) Scan(src interface{}) error {
	code, ok := src.(string)
	if !ok {
		return errors.New("currency code must be a string")
	}

	if !isCurrencyCode(code) {
		return errors.New("invalid currency code")
	}

	c.code = code
	return nil
}

func panicOnDifferentCurrency(first, second Amount) {
	if //noinspection GoBoolExpressions
	first.currency != second.currency && first.amountInMinor != 0 && second.amountInMinor != 0 {
		panic("Different currencies detected.")
	}
}

func resolveCurrency(first, second Currency) Currency {
	if first.isDefined() && !second.isDefined() {
		return first
	}

	if !first.isDefined() && second.isDefined() {
		return second
	}

	if first != second {
		panic(fmt.Errorf("same currency expected but got %s and %s", first, second))
	}

	return first
}
