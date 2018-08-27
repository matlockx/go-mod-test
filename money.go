package go-mod-test

import (
	"fmt"
)

var EUR = Currency{"EUR"}
var GBP = Currency{"GBP"}
var USD = Currency{"USD"}
var AUD = Currency{"AUD"}

type Amount struct {
	amountInMinor int
	currency      Currency
}

func Zero(currency Currency) Amount {
	return Amount{amountInMinor: 0, currency: currency}
}

func FromMinorInt64(amount int64, currency Currency) Amount {
	return Amount{int(amount), currency}
}

func FromMajorInt64(amount int64, currency Currency) Amount {
	return Amount{100 * int(amount), currency}
}

func FromMinor(amount int, currency Currency) Amount {
	return Amount{amount, currency}
}

func FromMajor(amount int, currency Currency) Amount {
	return Amount{100 * amount, currency}
}

func (m Amount) InMinorUnits() int {
	return m.amountInMinor
}

func (m Amount) InMinorUnits64() int64 {
	return int64(m.amountInMinor)
}

func (m Amount) InMajorUnits() float64 {
	return float64(m.amountInMinor) / 100.0
}

func (m Amount) Currency() Currency {
	return m.currency
}

func (m Amount) String() string {
	return fmt.Sprintf("%1.2f %s", m.InMajorUnits(), m.currency.code)
}

func (lhs Amount) Add(rhs Amount) Amount {
	return Amount{
		amountInMinor: lhs.amountInMinor + rhs.amountInMinor,
		currency:      resolveCurrency(lhs.currency, rhs.currency),
	}
}

func (lhs Amount) Multiply(factor int) Amount {
	return Amount{
		amountInMinor: lhs.amountInMinor * factor,
		currency:      lhs.currency,
	}
}

func (lhs Amount) LessThan(rhs Amount) bool {
	panicOnDifferentCurrency(lhs, rhs)
	return lhs.amountInMinor < rhs.amountInMinor
}
