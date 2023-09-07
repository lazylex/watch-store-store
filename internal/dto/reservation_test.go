package dto

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store/store/internal/dto/validators"
	"testing"
)

func TestReservationDTO_Validate(t *testing.T) {
	testCases := []struct {
		testName    string
		state       uint
		order       reservation.OrderNumber
		expectedErr error
	}{
		{
			testName:    "negative order number",
			state:       reservation.NewForCashRegister,
			order:       reservation.OrderNumber(-1),
			expectedErr: validators.ErrIncorrectOrder,
		},
		{
			testName:    "incorrect state",
			state:       200,
			order:       reservation.OrderNumber(1987),
			expectedErr: validators.ErrIncorrectState,
		},
		{
			testName:    "correct order for cash register",
			state:       reservation.NewForCashRegister,
			order:       reservation.OrderNumber(10),
			expectedErr: nil,
		},
		{
			testName:    "incorrect order for cash register",
			state:       reservation.NewForCashRegister,
			order:       reservation.OrderNumber(11),
			expectedErr: validators.ErrCashRegisterOrder,
		},
		{
			testName:    "correct order for local customer",
			state:       reservation.NewForLocalCustomer,
			order:       reservation.OrderNumber(11),
			expectedErr: nil,
		},
		{
			testName:    "incorrect order for local customer",
			state:       reservation.NewForLocalCustomer,
			order:       reservation.OrderNumber(10),
			expectedErr: validators.ErrOrderForInternetCustomer,
		},
		{
			testName:    "correct order for internet customer",
			state:       reservation.NewForInternetCustomer,
			order:       reservation.OrderNumber(11),
			expectedErr: nil,
		},
		{
			testName:    "incorrect order for internet customer",
			state:       reservation.NewForInternetCustomer,
			order:       reservation.OrderNumber(10),
			expectedErr: validators.ErrOrderForInternetCustomer,
		},
	}

	for _, tc := range testCases {
		r := &ReservationDTO{OrderNumber: tc.order, State: tc.state}
		t.Run(tc.testName, func(t *testing.T) {
			if !errors.Is(r.Validate(), tc.expectedErr) {
				t.Fail()
			}
		})
	}
}
