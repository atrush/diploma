package accrual

import (
	"context"
	"errors"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/provider/accrual"
	"github.com/atrush/diploma.git/services/order"
	"golang.org/x/time/rate"
	"log"
	"sync"
	"time"
)

type (
	AccrualService struct {
		svcOrders       order.OrderManager
		client          accrual.AccrualProvider
		producerTimeout time.Duration
		limiter         *rate.Limiter
	}

	LimiterChange struct {
		WaitSeconds int
		PerMinute   int
	}

	AccrualProcessor func(ctx context.Context, order model.Order, limitChan chan LimiterChange)

	Operator func(x float64) float64
)

var (
	batchSize = 10
)

func NewAccrualService(svc order.OrderManager, svcAcc accrual.AccrualProvider) *AccrualService {
	return &AccrualService{
		svcOrders:       svc,
		client:          svcAcc,
		producerTimeout: time.Duration(5) * time.Second,
		limiter:         rate.NewLimiter(rate.Limit(1), 1),
	}
}

func (s *AccrualService) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.producerTimeout)

	go func() {
		defer func() {
			log.Println("defer")
			ticker.Stop()
		}()
		for {
			log.Println("for")
			select {
			case <-ticker.C:
				batch, err := s.svcOrders.GetUnprocessedOrders(ctx, batchSize)
				log.Printf("start processing batch length:%v", len(batch))
				if err != nil {
					log.Printf("error getting new orders to process accruals, err:%s", err.Error())
					break
				}

				if len(batch) == 0 {
					log.Println("no new orders to process accruals")
					break
				}

				s.tick(ctx, batch, s.processAccrualForOrder)
				log.Printf("end processing batch length:%v", len(batch))
			case <-ctx.Done():
				log.Println("accrual worker context done")
				return
			}
		}
	}()
	return nil
}

//  tick gets accruals fo batch of orders, using function AccrualProcessor
//  if context is closed - tasks stops, waits all started, returns unprocessed to storage
//  if one of tasks returns change limit in chan - waits specified time, change limiter, return unprocessed to storage
func (s AccrualService) tick(ctx context.Context, batch []model.Order, process AccrualProcessor) {

	wg := sync.WaitGroup{}
	wg.Add(len(batch))

	// chanel listen change limit action
	chanLimit := make(chan LimiterChange, len(batch))
	defer close(chanLimit)

	lastProcessedIndex := 0

loop:
	for i, o := range batch {
		select {

		//  if received limiter change action - pauses, sets accrual limiter, stops tick
		case limit := <-chanLimit:
			perSec := float64(limit.PerMinute) / float64(60)
			log.Printf("changesleep %v, per second %v", limit.WaitSeconds, perSec)

			time.Sleep((time.Duration)(limit.WaitSeconds) * time.Second)
			s.limiter.SetLimit(rate.Limit(perSec))

			break loop

		//  on default - starts waits limiter, starts order process
		//  if context closed - limiter returns error, stops tick
		default:

			if err := s.limiter.Wait(ctx); err != nil {
				log.Printf("err %s: %s", o.Number, err.Error())

				break loop
			}

			lastProcessedIndex = i
			log.Printf("process %v, number: %v", lastProcessedIndex, o.Number)
			o := o
			go func() {
				log.Printf("go process %v, number: %v", lastProcessedIndex, o.Number)
				defer wg.Done()

				process(ctx, o, chanLimit)
			}()
		}
	}

	//  if exist not processed orders, clear change wait group count
	processedCount := lastProcessedIndex + 1
	if processedCount < len(batch) {
		returnOrders := batch[processedCount:]
		if err := s.svcOrders.ReturnNotUpdatedOrders(context.Background(), returnOrders); err != nil {
			log.Printf("error on returning unprocessed orders:%v", err.Error())
		}

		pass := (len(batch) - processedCount)
		log.Printf("not processed :%v", len(returnOrders))

		wg.Add(-1 * pass)
	}

	wg.Wait()
	log.Println("all ended")
}

//  processAccrualForOrder gets accrual from provider and updates order
//  if returns accrual - sends updates order status and accrual
//  if returns error ErrorAccrualLimitAchieved - writes LimiterChange in limiterChan
//  if returns error - sets to order status NEW back
func (s *AccrualService) processAccrualForOrder(ctx context.Context, order model.Order, limiterChan chan LimiterChange) {
	isProcessed := false
	ctxGet, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer func() {
		cancel()
		if !isProcessed {
			if err := s.svcOrders.UpdateStatus(context.Background(), order.ID, model.OrderStatusNew); err != nil {
				log.Printf("err deffer rollback order status for order number:%v, err:%s",
					order.Number, err.Error())
			}
		}
	}()

	accrualObj, err := s.client.Get(ctxGet, order.Number)
	if err != nil {
		log.Printf("error accrual get order number :%v, err:%s", order.Number, err.Error())

		if errors.Is(err, &model.ErrorAccrualLimitAchieved{}) {
			achieveErr, _ := err.(*model.ErrorAccrualLimitAchieved)
			limiterChan <- LimiterChange{
				WaitSeconds: achieveErr.WaitSeconds,
				PerMinute:   achieveErr.PerMinute,
			}
			return
		}
		return
	}

	if err := s.svcOrders.UpdateAccrual(context.TODO(), order, accrualObj); err != nil {
		log.Print("error accrual update")
		return
	}
	isProcessed = true
}
