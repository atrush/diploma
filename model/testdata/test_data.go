package testdata

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/pkg/luhn"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"
)

type TestData struct {
	Users    []model.User
	Orders   []model.Order
	Accruals []model.Accrual
}

const (
	testFile = "testdata.json"
)

func ReadTestData() (*TestData, error) {
	result := TestData{}

	data, err := ioutil.ReadFile(getPath())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GenOrders(count int, status model.OrderStatus) (uuid.UUID, []model.Order) {
	result := make([]model.Order, 0, count)
	userID := uuid.New()

	numbers := luhn.Generate(16, count)
	for i := 0; i < count; i++ {
		result = append(result, model.Order{
			ID:     uuid.New(),
			Number: numbers[i],
			UserID: userID,
			Status: status,
		})
	}
	return userID, result
}

func genNewDataFile() {
	testData := TestData{}

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 10; i++ {
		userID, userOrders := GenOrders(rand.Intn(5), model.OrderStatusNew)
		testData.Users = append(testData.Users, model.User{ID: userID, Login: fmt.Sprintf("a%v", i)})

		testData.Orders = append(testData.Orders, userOrders...)

	}

	for _, a := range testData.Orders {
		testData.Accruals = append(testData.Accruals, model.Accrual{
			Accrual: (rand.Intn(5) + 2) * 10000,
			Number:  a.Number,
			Status:  model.AccrualStatusProcessed,
		})
	}

	jsFixtures, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	WriteToFileMust(jsFixtures, "new_testdata.json")
}

func WriteToFileMust(data []byte, filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		log.Fatal(err.Error())
	}

	if err := writer.WriteByte('\n'); err != nil {
		log.Fatal("ошибка записи в файл: %w", err)
	}
}

// getFixturesDir returns current file directory.
func getPath() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	return path.Join(path.Dir(filePath), testFile)
}

//
//func TestWorker_tick(t *testing.T) {
//	genNewDataFile()
//	data := ReadTestDataMust()
//	log.Printf("%+v", data)
//}
