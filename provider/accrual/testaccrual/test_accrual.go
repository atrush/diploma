package testaccrual

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/model/testdata"
	"github.com/atrush/diploma.git/provider/accrual"
	"log"
	"math/rand"
	"time"
)

type TestAccrualProvider struct {
	Accruals []model.Accrual
}

var _ accrual.AccrualProvider = (*TestAccrualProvider)(nil)

func NewTestAccrualProvider() (*TestAccrualProvider, error) {

	d, err := testdata.ReadTestData()
	if err != nil {
		return nil, err
	}

	var b [8]byte
	_, err = crypto_rand.Read(b[:])
	if err != nil {
		log.Fatal(err.Error())
	}

	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	return &TestAccrualProvider{Accruals: d.Accruals}, nil
}

func (t *TestAccrualProvider) Get(ctx context.Context, number string) (model.Accrual, error) {
	for _, a := range t.Accruals {
		if a.Number == number {
			n := rand.Intn(4)
			time.Sleep(time.Duration(n) * time.Second)
			return a, nil
		}
	}
	return model.Accrual{}, nil
}
