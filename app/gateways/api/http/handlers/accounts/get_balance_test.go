package accountsroute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	type httpIO struct {
		r *http.Request
		w http.ResponseWriter
	}

	type testCase struct {
		name            string
		usecase         accountuc.Usecase
		httpIO          httpIO
		expectedStatus  int
		expectedPayload interface{}
	}

	var (
		testId = account.AccountId(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully get balance",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						fmt.Sprintf("/accounts/%s/balance", testId),
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {
					return 10, nil
				},
			},
			expectedStatus: 200,
			expectedPayload: map[string]interface{}{
				"account_id": testId.String(),
				"balance":    "R$0,10",
			},
		},
		{
			name: "successfully get 0 balance",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						fmt.Sprintf("/accounts/%s/balance", testId),
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {
					return 0, nil
				},
			},
			expectedStatus: 200,
			expectedPayload: map[string]interface{}{
				"account_id": testId.String(),
				"balance":    "R$0,00",
			},
		},
		{
			name: "fail to get balance for invalid id param",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts/1/balance",
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrInvalidPathParameter.Payload.Message},
		},
		{
			name: "fail to get balance inexistent account id",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						fmt.Sprintf("/accounts/%s/balance", testId),
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {
					return 0, account.ErrIdNotFound
				},
			},
			expectedStatus:  404,
			expectedPayload: map[string]interface{}{"title": responses.ErrAccountNotFound.Payload.Message},
		},
		{
			name: "fail due to usecase error",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						fmt.Sprintf("/accounts/%s/balance", testId),
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {
					return 0, accountuc.ErrRepository
				},
			},
			expectedStatus:  500,
			expectedPayload: map[string]interface{}{"title": responses.ErrInternalServerError.Payload.Message},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.usecase)

			h.GetBalance(tt.httpIO.w, tt.httpIO.r)

			recorder, ok := tt.httpIO.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			if statusCode := recorder.Code; statusCode != tt.expectedStatus {
				t.Errorf("got status code %d expected %d", statusCode, tt.expectedStatus)
			}

			var got map[string]interface{}
			err := json.NewDecoder(recorder.Body).Decode(&got)

			if err != nil {
				t.Fatalf("failed to decode response body: %s", err)
			}

			if !reflect.DeepEqual(got, tt.expectedPayload) {
				t.Fatalf("got response body: %s, expected %s", got, tt.expectedPayload)
			}
		})
	}

}
