package dto

import (
	"errors"
	"github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
	"testing"
)

func TestReservationDTO_Validate(t *testing.T) {
	testCases := []struct {
		testName    string
		state       uint
		order       reservation.OrderNumber
		products    []ProductDTO
		expectedErr error
	}{
		{
			testName:    "negative order number",
			state:       reservation.NewForCashRegister,
			order:       reservation.OrderNumber(-1),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: validators.ErrIncorrectOrder,
		},
		{
			testName:    "incorrect state",
			state:       200,
			order:       reservation.OrderNumber(1987),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: validators.ErrIncorrectState,
		},
		{
			testName:    "correct order for cash register",
			state:       reservation.NewForCashRegister,
			order:       reservation.OrderNumber(reservation.MaxCashRegisterNumber),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: nil,
		},
		{
			testName:    "incorrect order for cash register",
			state:       reservation.NewForCashRegister,
			order:       reservation.OrderNumber(reservation.MaxCashRegisterNumber + 1),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: validators.ErrCashRegisterOrder,
		},
		{
			testName:    "correct order for local customer",
			state:       reservation.NewForLocalCustomer,
			order:       reservation.OrderNumber(reservation.MaxCashRegisterNumber + 1),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: nil,
		},
		{
			testName:    "incorrect order for local customer",
			state:       reservation.NewForLocalCustomer,
			order:       reservation.OrderNumber(reservation.MaxCashRegisterNumber),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: validators.ErrOrderForInternetCustomer,
		},
		{
			testName:    "correct order for internet customer",
			state:       reservation.NewForInternetCustomer,
			order:       reservation.OrderNumber(reservation.MaxCashRegisterNumber + 1),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: nil,
		},
		{
			testName:    "incorrect order for internet customer",
			state:       reservation.NewForInternetCustomer,
			order:       reservation.OrderNumber(reservation.MaxCashRegisterNumber),
			products:    []ProductDTO{{Article: "ca-09.1000", Price: 4660, Amount: 5}},
			expectedErr: validators.ErrOrderForInternetCustomer,
		},
		{
			testName:    "empty products",
			state:       reservation.NewForCashRegister,
			order:       reservation.OrderNumber(reservation.MaxCashRegisterNumber),
			products:    []ProductDTO{},
			expectedErr: validators.ErrNoProductsInReservation,
		},
	}

	for _, tc := range testCases {
		r := &ReservationDTO{OrderNumber: tc.order, State: tc.state, Products: tc.products}
		t.Run(tc.testName, func(t *testing.T) {
			if !errors.Is(r.Validate(), tc.expectedErr) {
				t.Fail()
			}
		})
	}

	t.Run("duplicate product", func(t *testing.T) {
		r := &ReservationDTO{OrderNumber: 466, State: reservation.NewForInternetCustomer,
			Products: []ProductDTO{
				{
					Article: "ca-09.1000",
					Price:   4660,
					Amount:  5,
				},
				{
					Article: "ca-12",
					Price:   46960,
					Amount:  1,
				},
				{
					Article: "ca-09.1000",
					Price:   4660,
					Amount:  9,
				},
			}}
		if !errors.Is(r.Validate(), validators.ErrDuplicateProductsInReservation) {
			t.Fail()
		}
	})

	t.Run("error product", func(t *testing.T) {
		r := &ReservationDTO{OrderNumber: 466, State: reservation.NewForInternetCustomer,
			Products: []ProductDTO{
				{
					Article: "ca-09.1900",
					Price:   4660,
					Amount:  5,
				},
			}}
		if !errors.Is(r.Validate(), validators.ErrIncorrectArticle) {
			t.Fail()
		}
	})

}
