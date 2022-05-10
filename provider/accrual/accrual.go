package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	model_api "github.com/atrush/diploma.git/api/model"
	"github.com/atrush/diploma.git/model"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Accrual struct {
	serviceURL string
}

var _ AccrualProvider = (*Accrual)(nil)

//  newAccrualWithHTTP returns new accrual, accepts http client interface for testing
func NewAccrual(url string) (*Accrual, error) {
	return &Accrual{
		serviceURL: url,
	}, nil
}

// RequestAccrual
func (a *Accrual) Get(ctx context.Context, number string) (model.Accrual, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", a.serviceURL, number),
		nil,
	)
	if err != nil {
		log.Printf("client error %v", err.Error())
		return model.Accrual{}, fmt.Errorf("error accrual request %w", err)
	}
	client := http.Client{}
	r, err := client.Do(request)

	if err != nil {
		log.Printf("get error %v", err.Error())
		return model.Accrual{}, fmt.Errorf("error accrual request %w", err)
	}
	log.Printf("resp%v status", r.StatusCode)
	// 200 parse and check response
	if r.StatusCode == http.StatusOK {
		var respObj model_api.AccrualResponse

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := decoder.Decode(&respObj); err != nil {
			log.Printf("accrual decode err %v", err.Error())
			return model.Accrual{}, fmt.Errorf("error accrual request: error decode json response:%w", err)
		}

		//  convert to canonical and checks
		accrual, err := respObj.ToCanonical()
		if err != nil {
			return model.Accrual{}, fmt.Errorf("error accrual request: %w", err)
		}

		if accrual.Number != number {
			return model.Accrual{},
				fmt.Errorf("error accrual request: response number %v does not match order bumber %v",
					accrual.Number, number)
		}
		log.Printf("accrual %+v parsed", accrual)
		return accrual, nil

	}

	if r.StatusCode == http.StatusInternalServerError {
		return model.Accrual{}, errors.New("error accrual fetch - 500 response")
	}

	if r.StatusCode == http.StatusTooManyRequests {
		strWait := r.Header.Values("retry-after")[0]
		waitSec, err := strconv.Atoi(strWait)
		if err != nil {
			return model.Accrual{}, errors.New("error accrual fetch - 500 response")
		}

		resBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error processing 429 response from accrual service")
			return model.Accrual{}, fmt.Errorf("error processing 429 response from accrual service: %w", err)
		}

		defer r.Body.Close()

		re := regexp.MustCompile("[0-9]+")
		strPerMinute := re.FindAllString(string(resBody), 1)[0]
		perMinute, err := strconv.Atoi(strPerMinute)
		if err != nil {
			return model.Accrual{}, errors.New("error accrual fetch - 500 response")
		}

		return model.Accrual{}, &model.ErrorAccrualLimitAchieved{
			WaitSeconds: waitSec,
			PerMinute:   perMinute,
		}
	}

	//- `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
	//- `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
	//- `PROCESSING` — расчёт начисления в процессе;
	//- `PROCESSED` — расчёт начисления окончен;

	// 429 - too many requests, pause
	// if 500 or unexpected status return error
	return model.Accrual{}, fmt.Errorf("error accrual request: response status code:%v", r.StatusCode)
}
