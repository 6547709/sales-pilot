package view

import "strings"

// Mode 内容视图：销售视图（默认）或客户视图
type Mode string

const (
	Sales    Mode = "sales"
	Customer Mode = "customer"
)

func Parse(raw string) Mode {
	if strings.EqualFold(strings.TrimSpace(raw), string(Customer)) {
		return Customer
	}
	return Sales
}

func (m Mode) IsCustomer() bool {
	return m == Customer
}
