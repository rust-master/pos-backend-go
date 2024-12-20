package models

type LoginAdmin struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Products struct {
	ProductCode  int64   `json:"product_code"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	UnitPrice    float64 `json:"unitprice"`
	UnitsInStock int64   `json:"unitsinstock"`
}

type Cutomers struct {
	CustomerID  int64  `json:"customerID"`
	Name        string `json:"name"`
	MobilePhone string `json:"mobilePhone"`
}

type Orders struct {
	OrderID     int64  `json:"orderID"`
	CustomerID  int64  `json:"customerID"`
	Date        uint64 `json:"date"`
	OrderStatus string `json:"orderStatus"`
}

type Invoices struct {
	InvoiceID     int64  `json:"invoiceID"`
	OrderID       int64  `json:"orderID"`
	Date          uint64 `json:"date"`
	InvoiceStatus string `json:"invoiceStatus"`
}

type InvoiceLines struct {
	InvoiceID   int64  `json:"invoiceID"`
	ProductCode int64  `json:"productCode"`
	UnitPrice   uint64 `json:"unitPrice"`
}
