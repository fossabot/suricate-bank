package accountspg

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestGetById(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type testCase struct {
		name            string
		runBefore       func() error
		idArg           account.AccountId
		expectedAccount account.Account
		err             error
	}

	var (
		testIdCredit10initial0 = account.AccountId(uuid.New())
		testCpf                = cpf.Random()
	)

	testCases := []testCase{
		{
			name: "successfully get an account",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIdCredit10initial0,
					testCpf.Value(),
					0,
				)
			},
			idArg: testIdCredit10initial0,
			expectedAccount: account.Account{
				Id:        testIdCredit10initial0,
				Name:      "nice name",
				Cpf:       testCpf,
				Secret:    testHash,
				CreatedAt: testTime,
			},
		},
		{
			name:            "fail to get an inexixtent account",
			idArg:           account.AccountId(uuid.New()),
			expectedAccount: account.Account{},
			err:             ErrQuery,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				err := tt.runBefore()

				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			gotAccount, err := testRepo.GetById(testContext, tt.idArg)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if !reflect.DeepEqual(gotAccount, tt.expectedAccount) {
				t.Fatalf("got %v expected %v", gotAccount, tt.expectedAccount)
			}
		})
	}
}
