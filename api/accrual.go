package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/services/httpclient"
	"github.com/atrush/diploma.git/services/order"
	"net/http"
)

type Accrual struct {
	svc        order.OrderManager
	serviceURL string
	httpClient httpclient.HTTPClient
}

//  NewAccrual returns new accrual,
func NewAccrual(svc order.OrderManager, url string) (*Accrual, error) {
	return newAccrualWithHTTP(svc, url, &http.Client{})
}

//  newAccrualWithHTTP returns new accrual, accepts http client interface for testing
func newAccrualWithHTTP(svc order.OrderManager, url string, client httpclient.HTTPClient) (*Accrual, error) {
	return &Accrual{
		svc:        svc,
		serviceURL: url,
		httpClient: client, // interface for test mocking
	}, nil
}

func (a *Accrual) ProcessOrder(ctx context.Context, order model.Order) (err error) {
	//  change order status on start processing
	if err = a.svc.UpdateStatus(ctx, order.ID, model.OrderStatusProcessing); err != nil {
		err = fmt.Errorf("error process order accrual: error update order status:%w", err)
		return err
	}

	//  if was error change status back to new
	defer func() {
		if err != nil {
			errDef := a.svc.UpdateStatus(ctx, order.ID, model.OrderStatusNew)
			if errDef != nil {
				err = fmt.Errorf(err.Error()+": error deff: %w", errDef)
			}
		}
	}()

	accrual, err := a.GetAccrualRequest(ctx, order)
	if err != nil {
		err = fmt.Errorf("error process order accrual:%w", err)
		return err
	}

	if err := a.svc.UpdateAccrual(ctx, order, accrual); err != nil {
		err = fmt.Errorf("error process order accrual:%w", err)
		return err
	}
	return nil
}

// RequestAccrual
func (a *Accrual) GetAccrualRequest(ctx context.Context, order model.Order) (model.Accrual, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v/api/orders/%v", a.serviceURL, order.Number),
		nil,
	)
	if err != nil {
		return model.Accrual{}, fmt.Errorf("error accrual request %w", err)
	}

	r, err := a.httpClient.Do(request)

	// 200 parse and check response
	if r.StatusCode == http.StatusOK {
		var respObj AccrualResponse

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := decoder.Decode(&respObj); err != nil {
			return model.Accrual{}, fmt.Errorf("error accrual request: error decode json response:%w", err)
		}

		//  convert to canonical and checks
		accrual, err := respObj.ToCanonical()
		if err != nil {
			return model.Accrual{}, fmt.Errorf("error accrual request: %w", err)
		}

		if accrual.Number != order.Number {
			return model.Accrual{},
				fmt.Errorf("error accrual request: response number %v does not match order bumber %v",
					accrual.Number, order.Number)
		}

		return accrual, nil

	}

	//- `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
	//- `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
	//- `PROCESSING` — расчёт начисления в процессе;
	//- `PROCESSED` — расчёт начисления окончен;

	// 429 - too many requests, pause
	// if 500 or unexpected status return error
	return model.Accrual{}, fmt.Errorf("error accrual request: response status code:%v", r.StatusCode)
}
