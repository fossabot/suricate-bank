package transferuc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository transfer.Repository
		want       []transfer.Transfer
		err        error
	}

	var testUUID1 = uuid.New()
	var testUUID2 = uuid.New()
	var testUUID3 = uuid.New()
	var testUUID4 = uuid.New()

	testCases := []testCase{
		{
			name: "successfully fetch 1 transfer",
			repository: transfer.MockRepository{
				OnFetch: func(ctx context.Context) ([]transfer.Transfer, error) {

					return []transfer.Transfer{
						{
							Id: transfer.TransferId(testUUID1),
						},
					}, nil
				},
			},
			want: []transfer.Transfer{
				{
					Id: transfer.TransferId(testUUID1),
				},
			},
		},
		{
			name: "successfully fetch 4 transfers",
			repository: transfer.MockRepository{
				OnFetch: func(ctx context.Context) ([]transfer.Transfer, error) {

					return []transfer.Transfer{
						{
							Id: transfer.TransferId(testUUID1),
						},
						{
							Id: transfer.TransferId(testUUID2),
						},
						{
							Id: transfer.TransferId(testUUID3),
						},
						{
							Id: transfer.TransferId(testUUID4),
						},
					}, nil
				},
			},
			want: []transfer.Transfer{
				{
					Id: transfer.TransferId(testUUID1),
				},
				{
					Id: transfer.TransferId(testUUID2),
				},
				{
					Id: transfer.TransferId(testUUID3),
				},
				{
					Id: transfer.TransferId(testUUID4),
				},
			},
		},
		{
			name: "fetch zero transfers",
			repository: transfer.MockRepository{
				OnFetch: func(ctx context.Context) ([]transfer.Transfer, error) {

					return nil, nil
				},
			},
			want: nil,
		},
		{
			name: "repository throws error",
			repository: transfer.MockRepository{
				OnFetch: func(ctx context.Context) ([]transfer.Transfer, error) {

					return []transfer.Transfer{}, ErrRepository
				},
			},
			want: []transfer.Transfer{},
			err:  ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository, MockCrediter{}, MockDebiter{}}

			transfersList, err := uc.Fetch(context.Background())

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(transfersList, tt.want) {
				t.Errorf("got %v expected %v", transfersList, tt.want)
			}
		})
	}
}
