package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type PaymentResult struct {
	TransactionID string
	Amount        float64
	Method        string
	Success       bool
}

type PaymentProcessor interface {
	Name() string
	Type() string
	Charge(amount float64, reference string) (*PaymentResult, error)
	Refund(transactionID string) error
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type ConsoleLogger struct {
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}
func (l *ConsoleLogger) Info(msg string)  { fmt.Println(capitalize(msg)) }
func (l *ConsoleLogger) Warn(msg string)  { fmt.Println("WARN: " + capitalize(msg)) }
func (l *ConsoleLogger) Error(msg string) { fmt.Println(capitalize(msg)) }

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

type CreditCard struct {
	CardNumber string // last 4 digits only for display
	HolderName string
}

func NewCreditCard(number, holder string) *CreditCard {
	if len(number) < 4 {
		return nil
	}
	lastFour := number[len(number)-4:]
	return &CreditCard{
		CardNumber: lastFour,
		HolderName: holder,
	}
}

func (c *CreditCard) Name() string {
	return fmt.Sprintf("CreditCard (**** %s)", c.CardNumber)
}

func (c *CreditCard) Type() string { return "CreditCard" }

func (c *CreditCard) Charge(amount float64, reference string) (*PaymentResult, error) {
	const limit = 100000.0
	if amount <= 0 {
		return nil, errors.New("charge amount must be greater than zero")
	}
	if amount > limit {
		return nil, fmt.Errorf("amount %.f exceeds limit of %.f", amount, limit)
	}
	return &PaymentResult{
		TransactionID: fmt.Sprintf("txn CC-%d", rand.Intn(900000)+100000),
		Amount:        amount,
		Method:        c.Name(),
		Success:       true,
	}, nil
}

func (c *CreditCard) Refund(transactionID string) error { return nil }

type JazzCash struct {
	PhoneNumber string
}

func NewJazzCash(phone string) *JazzCash {
	return &JazzCash{PhoneNumber: phone}
}

func (j *JazzCash) Name() string {
	return fmt.Sprintf("Jazzcash %s", j.PhoneNumber)
}

func (j *JazzCash) Type() string { return "JazzCash" }

func (j *JazzCash) Charge(amount float64, reference string) (*PaymentResult, error) {
	limit := 25000.0
	if amount <= 0 {
		return nil, errors.New("charge amount must be greater than zero")
	}
	if amount > limit {
		return nil, fmt.Errorf("amount %.f exceeds wallet limit of 25000", amount)
	}
	return &PaymentResult{
		TransactionID: fmt.Sprintf("txn WT-%d", rand.Intn(900000)+100000),
		Amount:        amount,
		Method:        j.Name(),
		Success:       true,
	}, nil
}

func (j *JazzCash) Refund(transactionID string) error { return nil }

type BankTransfer struct {
	AccountNumber string
	BankName      string
}

func NewBankTransfer(account, bank string) *BankTransfer {
	return &BankTransfer{AccountNumber: account, BankName: bank}
}

func (b *BankTransfer) Name() string {
	return fmt.Sprintf("Account Number %s", b.AccountNumber)
}

func (b *BankTransfer) Type() string { return "BankTransfer" }

func (b *BankTransfer) Charge(amount float64, reference string) (*PaymentResult, error) {
	if amount <= 0 {
		return nil, errors.New("charge amount must be greater than zero")
	}
	return &PaymentResult{
		TransactionID: fmt.Sprintf("txn BT-%d", rand.Intn(900000)+100000),
		Amount:        amount,
		Method:        b.Name(),
		Success:       true,
	}, nil
}

func (b *BankTransfer) Refund(transactionID string) error { return nil }

type CheckoutService struct {
	processor PaymentProcessor
	logger    Logger
}

func NewCheckoutService(processor PaymentProcessor, logger Logger) *CheckoutService {
	return &CheckoutService{
		processor: processor,
		logger:    logger,
	}
}

func (c *CheckoutService) Checkout(orderID string, amount float64) (*PaymentResult, error) {
	c.logger.Info(fmt.Sprintf("charging PKR %.0f via %s for order %s", amount, c.processor.Name(), orderID))
	result, err := c.processor.Charge(amount, orderID)
	if err != nil {
		c.logger.Error(fmt.Sprintf("payment failed: %s: %s", c.processor.Type(), err.Error()))
		return nil, err
	}
	c.logger.Info(fmt.Sprintf("payment successful: txn %s", result.TransactionID))
	return result, nil
}

func main() {
	//CREDIT CARD RUN
	logger := &ConsoleLogger{}
	cardObj := NewCreditCard("1234567890124242", "Ali Ahmed")
	cardCheckout := NewCheckoutService(cardObj, logger)
	cardCheckout.Checkout("ORD-001", 5000)

	fmt.Println()
	//JAZZCASH RUN
	jazzObj := NewJazzCash("03001234567")
	jazzCashCheckout := NewCheckoutService(jazzObj, logger)
	jazzCashCheckout.Checkout("ORD-002", 30000)

	fmt.Println()
	// BANK TRANSFER RUN
	bankObj := NewBankTransfer("PK12BANK00001111", "HBL")
	bankCheckout := NewCheckoutService(bankObj, logger)
	bankCheckout.Checkout("ORD-003", 30000)
}
