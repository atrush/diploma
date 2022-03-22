package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/services/httpclient"
	mkHTTP "github.com/atrush/diploma.git/services/httpclient/mock"
	"github.com/atrush/diploma.git/services/order"
	"github.com/atrush/diploma.git/storage"
	mkSt "github.com/atrush/diploma.git/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

type accrualTest struct {
	name       string
	st         storage.Storage
	client     httpclient.HTTPClient
	expAccrual int
	expStatus  model.AccrualStatus
}

func TestAccrual_ProcessOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrder := model.Order{
		ID:     uuid.New(),
		Number: "5336610979490109",
		Status: model.OrderStatusNew,
	}

	tests := []struct {
		name            string
		st              storage.Storage
		client          httpclient.HTTPClient
		expAccrual      int
		expStatus       model.AccrualStatus
		expErr          bool
		expErrType      bool
		expErrTypeCheck error
	}{
		{
			name:   "processed with accrual 500",
			client: accrualProcessed(ctrl, `{"order": "5336610979490109","status": "PROCESSED","accrual": 500}`),
			st:     storageOrderUpdateAccrual(ctrl, model.OrderStatusProcessed, 500*model.MoneyAccuracy),
			expErr: false,
		},
		{
			name:   "processed with no accrual",
			client: accrualProcessed(ctrl, `{"order": "5336610979490109","status": "PROCESSED"}`),
			st:     storageOrderUpdateAccrual(ctrl, model.OrderStatusProcessed, 0),
			expErr: false,
		},
		{
			name:   "registered - no accrual, processing",
			client: accrualProcessed(ctrl, `{"order": "5336610979490109","status": "REGISTERED"}`),
			st:     storageOrderUpdateAccrual(ctrl, model.OrderStatusProcessing, 0),
			expErr: false,
		},
		{
			name:   "invalid - no accrual, invalid",
			client: accrualProcessed(ctrl, `{"order": "5336610979490109","status": "INVALID"}`),
			st:     storageOrderUpdateAccrual(ctrl, model.OrderStatusInvalid, 0),
			expErr: false,
		},
		{
			name:   "processing - no accrual, processing",
			client: accrualProcessed(ctrl, `{"order": "5336610979490109","status": "PROCESSING"}`),
			st:     storageOrderUpdateAccrual(ctrl, model.OrderStatusProcessing, 0),
			expErr: false,
		},
		{
			name:   "storage error - update status to new",
			client: accrualProcessed(ctrl, `{"order": "5336610979490109","status": "PROCESSING"}`),
			st:     storageOrderUpdateAccrualError(ctrl, model.OrderStatusProcessing, 0),
			expErr: true,
		},

		{
			name:   "storage error - update status to new",
			client: accrualProcessed(ctrl, `{"order": "5336610979490109","status": "PROCESSING"}`),
			st:     storageOrderUpdateAccrualError(ctrl, model.OrderStatusProcessing, 0),
			expErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcOrder, err := order.NewOrder(tt.st)
			require.NoError(t, err)

			accr, err := newAccrualWithHTTP(svcOrder, "", tt.client)
			require.NoError(t, err)

			err = accr.ProcessOrder(context.Background(), mockOrder)
			if tt.expErr {
				require.Error(t, err)
				if tt.expErrType {
					require.True(t, errors.Is(err, tt.expErrTypeCheck))
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func storageOrderUpdateAccrual(ctrl *gomock.Controller, status model.OrderStatus, accrual int) *mkSt.MockStorage {
	mockOrderRepo := mkSt.NewMockOrderRepository(ctrl)
	mockOrderRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), model.OrderStatusProcessing).Return(nil)
	mockOrderRepo.EXPECT().UpdateAccrual(gomock.Any(), gomock.Any(), status, accrual).Return(nil)

	mockStorage := mkSt.NewMockStorage(ctrl)
	mockStorage.EXPECT().Order().Return(mockOrderRepo)
	mockStorage.EXPECT().Order().Return(mockOrderRepo)
	return mockStorage
}
func storageOrderUpdateAccrualError(ctrl *gomock.Controller, status model.OrderStatus, accrual int) *mkSt.MockStorage {
	mockOrderRepo := mkSt.NewMockOrderRepository(ctrl)
	mockOrderRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), model.OrderStatusProcessing).Return(nil)
	mockOrderRepo.EXPECT().UpdateAccrual(gomock.Any(), gomock.Any(), status, accrual).Return(fmt.Errorf("internal error"))
	mockOrderRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), model.OrderStatusNew).Return(nil)

	mockStorage := mkSt.NewMockStorage(ctrl)
	mockStorage.EXPECT().Order().Return(mockOrderRepo)
	mockStorage.EXPECT().Order().Return(mockOrderRepo)
	mockStorage.EXPECT().Order().Return(mockOrderRepo)
	return mockStorage
}

func accrualProcessed(ctrl *gomock.Controller, jsAccrual string) *mkHTTP.MockHTTPClient {
	mock := mkHTTP.NewMockHTTPClient(ctrl)
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsAccrual)))

	mock.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)
	return mock
}
