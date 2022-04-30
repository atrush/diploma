package handler

import (
	"fmt"
	"github.com/atrush/diploma.git/model"
	mk "github.com/atrush/diploma.git/services/order/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

var (
	userID     = uuid.New()
	userToken  = GenJWTAuthToken(userID)
	mockOrders = []model.Order{
		{
			ID:         uuid.New(),
			UserID:     userID,
			Number:     "5336610979490109",
			Accrual:    50060,
			Status:     model.OrderStatusProcessing,
			UploadedAt: time.Date(2021, 10, 1, 11, 11, 1, 0, time.UTC),
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			Number:     "6221273347496047",
			Status:     model.OrderStatusNew,
			Accrual:    50000,
			UploadedAt: time.Date(2021, 10, 3, 11, 11, 1, 0, time.UTC),
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			Number:     "4977253909241263",
			Status:     model.OrderStatusRegistered,
			UploadedAt: time.Date(2021, 10, 4, 11, 11, 1, 0, time.UTC),
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			Number:     "676194491210",
			Status:     model.OrderStatusInvalid,
			UploadedAt: time.Date(2021, 10, 5, 11, 11, 1, 0, time.UTC),
		},
	}
)

//  TestHandler_Login tests user register handler
func TestHandler_OrderGetListForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlOrderGet := "/api/user/orders"

	tests := []TestRoute{
		{
			name:            "return 200 if orders is exist for user",
			method:          http.MethodGet,
			url:             urlOrderGet,
			svcOrder:        mockOrderGetListForUser(ctrl),
			headers:         map[string]string{"Authorization": "Bearer " + userToken},
			expectedHeaders: map[string]string{"Content-Type": "application/json"},
			expectedCode:    200,
		},
		{
			name:         "return 204 if orders not exist for user",
			method:       http.MethodGet,
			url:          urlOrderGet,
			svcOrder:     mockOrderGetEmptyForUser(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			expectedCode: 204,
		},
		{
			name:         "return 401 user not authenticated",
			method:       http.MethodGet,
			url:          urlOrderGet,
			svcOrder:     mockOrderEmpty(ctrl),
			expectedCode: 401,
		},
		{
			name:         "return 500 if internal error",
			method:       http.MethodGet,
			url:          urlOrderGet,
			svcOrder:     mockOrderGetServerErrForUser(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}

}

//  TestHandler_Login tests user register handler
func TestHandler_OrderAddToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlOrderAdd := "/api/user/orders"

	tests := []TestRoute{
		{
			name:         "return 200 if order is exist for user",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderExistForUser(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			body:         mockOrders[0].Number,
			expectedCode: 200,
		},
		{
			name:         "return 202 number of order accepted",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderSaved(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			body:         mockOrders[0].Number,
			expectedCode: 202,
		},
		{
			name:         "return 400 wrong request if empty body",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderEmpty(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			expectedCode: 400,
		},
		{
			name:         "return 400 wrong request if empty body, if contains no-numbers",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderEmpty(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			body:         "dsfds32324324432",
			expectedCode: 400,
		},
		{
			name:         "return 401 user not authenticated",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderEmpty(ctrl),
			body:         mockOrders[0].Number,
			expectedCode: 401,
		},
		{
			name:         "return 409 number of order was uploaded by another user",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderExistForAnotherUser(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			body:         mockOrders[0].Number,
			expectedCode: 409,
		},
		{
			name:         "return 422 wrong order number",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderEmpty(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			body:         "11111111111",
			expectedCode: 422,
		},

		{
			name:         "return 500 server error",
			method:       http.MethodPost,
			url:          urlOrderAdd,
			svcOrder:     mockOrderServerError(ctrl),
			headers:      map[string]string{"Authorization": "Bearer " + userToken},
			body:         mockOrders[0].Number,
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}

}

/*  mocks  */
func mockOrderEmpty(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	return orderMock
}

/*  Get mocks  */
func mockOrderGetListForUser(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	orderMock.EXPECT().GetForUser(gomock.Any(), userID).Return(mockOrders, nil)
	return orderMock
}
func mockOrderGetEmptyForUser(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	orderMock.EXPECT().GetForUser(gomock.Any(), userID).Return(nil, nil)
	return orderMock
}
func mockOrderGetServerErrForUser(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	orderMock.EXPECT().GetForUser(gomock.Any(), userID).Return(nil, fmt.Errorf("internal errror"))
	return orderMock
}

/*  Create mocks  */
func mockOrderExistForUser(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	orderMock.EXPECT().AddToUser(gomock.Any(), gomock.Any(), userID).Return(model.Order{}, model.ErrorOrderExist)
	return orderMock
}
func mockOrderExistForAnotherUser(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	orderMock.EXPECT().AddToUser(gomock.Any(), gomock.Any(), userID).Return(model.Order{}, model.ErrorOrderExistAnotheUser)
	return orderMock
}
func mockOrderSaved(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	orderMock.EXPECT().AddToUser(gomock.Any(), gomock.Any(), userID).Return(mockOrders[0], nil)
	return orderMock
}
func mockOrderServerError(ctrl *gomock.Controller) *mk.MockOrderManager {
	orderMock := mk.NewMockOrderManager(ctrl)
	orderMock.EXPECT().AddToUser(gomock.Any(), gomock.Any(), userID).Return(model.Order{}, fmt.Errorf("internal errror"))
	return orderMock
}
