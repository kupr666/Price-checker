package service

import(
	"context"
	"testing"
	"price_checker/internal/core/domains"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go.uber.org/zap"
)

type MockRepository struct {
	// mock.Mock gives us 2 mothods On() & Called()
	mock.Mock
}

func (m *MockRepository) Add(ctx context.Context, item domains.Item) (domains.Item, error) {
	args := m.Called(ctx, item)
	return args.Get(0).(domains.Item), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetAll(ctx context.Context) ([]domains.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domains.Item), args.Error(1)
}

func (m *MockRepository) UpdatePrice(ctx context.Context, id int64, newPrice float64) error {
	args := m.Called(ctx, id, newPrice)
	return args.Error(0)
}

type MockScraper struct {
	mock.Mock
}

func (m *MockScraper) FetchCurrentPrice(ctx context.Context, itemURL string) (float64, error) {
	args := m.Called(ctx, itemURL)
	return args.Get(0).(float64), args.Error(1)
}

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) Notify(message string) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestCreateItem_Success(t *testing.T) {

	// ARRANGE 
	repo := new(MockRepository)
	scraper := new(MockScraper)
	tgNotifier := new(MockNotifier)
	logger := zap.NewNop() // fake logger (without any logs)

	svc := NewPriceService(repo, scraper, logger, tgNotifier)

	ctx := context.Background()
	inputItem := domains.Item{
		URL: "https://future-phone.ru/catalog/smartfony_i_gadzhety/smartfony/apple/iphone_se_2022/iphone_se_2022_64gb_starlight" ,
		TargetPrice: 30000,	
	} 

	expectedItem := domains.Item{
		ID: 1,
		URL: "https://future-phone.ru/catalog/smartfony_i_gadzhety/smartfony/apple/iphone_se_2022/iphone_se_2022_64gb_starlight,", 
		CurrentPrice: 26900,
		TargetPrice: 25000,
	}

	scraper.On("FetchCurrentPrice", ctx, inputItem.URL).Return(26900.0, nil)
	repo.On("Add", ctx, mock.AnythingOfType("domains.Item")).Return(expectedItem, nil)

	// ACT
	// call of function which we are testing
	result, err := svc.CreateItem(ctx, inputItem)
	
	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, expectedItem, result)

	repo.AssertExpectations(t)
	scraper.AssertExpectations(t)
}