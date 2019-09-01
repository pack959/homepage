package main

import (
	"errors"
	"path"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

type StripeProcessor struct {
	publicURL string
}

func NewStripeProcessor(publicURL, stripeKey string) *StripeProcessor {
	stripe.Key = stripeKey

	return &StripeProcessor{
		publicURL: publicURL,
	}
}

type Session struct {
	checkoutSession *stripe.CheckoutSession
	checkoutParams  *stripe.CheckoutSessionParams
	processor       *StripeProcessor
}

func (p StripeProcessor) CreateSession(successPath, cancelPath string) *Session {
	successURL := path.Join(p.publicURL, successPath)
	cancelURL := path.Join(p.publicURL, cancelPath)

	return &Session{
		processor: &p,
		checkoutParams: &stripe.CheckoutSessionParams{
			PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
			LineItems:          []*stripe.CheckoutSessionLineItemParams{},
			SuccessURL:         &successURL,
			CancelURL:          &cancelURL,
		},
		checkoutSession: nil,
	}
}

func (s *Session) AddItem(name, description string, amount, quantity int64) {
	item := &stripe.CheckoutSessionLineItemParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Quantity:    stripe.Int64(quantity),
	}

	s.checkoutParams.LineItems = append(s.checkoutParams.LineItems, item)
}

func (s *Session) Start() error {
	if s.checkoutParams == nil {
		return errors.New("missing stripe checkout params")
	}

	if s.checkoutParams.SuccessURL == nil {
		return errors.New("stripe session without successURL")
	}

	if s.checkoutParams.CancelURL == nil {
		return errors.New("stripe session without cancelURL")
	}

	if len(s.checkoutParams.LineItems) < 1 {
		return errors.New("stripe session without items")
	}

	session, err := session.New(s.checkoutParams)
	if err != nil {
		return err
	}

	s.checkoutSession = session
	return nil
}

func (s *Session) GetID() (string, error) {
	if s.checkoutSession == nil {
		return "", errors.New("stripe session id not found, try to start session")
	}
	return s.checkoutSession.ID, nil
}
