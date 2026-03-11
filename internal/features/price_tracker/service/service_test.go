package service

import (
	"fmt"
	"context"
	"price_checker/internal/core/domains"
	// "price_checker/internal/features/price_tracker/scraper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go.uber.org/zap"
)

const testURL = "https://future-phone.ru/catalog/smartfony_i_gadzhety/smartfony/apple/iphone_se_2022/iphone_se_2022_64gb_starlight"

type MockRepository struct {
	// mock.Mock gives us 2 methods On() & Called()
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
	notifier := new(MockNotifier)
	logger := zap.NewNop() // fake logger (without any logs)

	svc := NewPriceService(repo, scraper, logger, notifier)

	ctx := context.Background()
	inputItem := domains.Item{
		URL: testURL,
		TargetPrice: 30000,	
	} 

	expectedItem := domains.Item{
		ID: 1,
		URL: testURL, 
		CurrentPrice: 26900,
		TargetPrice: 30000,
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

func TestCreateItem_ScraperFailed(t *testing.T) {

	// ARRANGE
	repo := new(MockRepository)
	scraper := new(MockScraper)
	notifier := new(MockNotifier)
	logger := zap.NewNop()

	svc := NewPriceService(repo, scraper, logger, notifier)

	ctx := context.Background()
	inputItem := domains.Item {
		URL: testURL,
		TargetPrice: 26900,
	}

	scraper.On("FetchCurrentPrice", ctx, inputItem.URL).Return(0.0, fmt.Errorf("site unavailable"))

	// ACT
	result, err := svc.CreateItem(ctx, inputItem)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, domains.Item{}, result)
	repo.AssertNotCalled(t, "Add")
	scraper.AssertExpectations(t)
}

func TestCheckAllPrices_TargetReached(t *testing.T) {
	
	// ARRANGE
	repo := new(MockRepository)
	scraper := new(MockScraper)
	notifier := new(MockNotifier)
	logger := zap.NewNop()

	svc := NewPriceService(repo, scraper, logger, notifier)

	ctx := context.Background()
	inputItems := []domains.Item{
		{ID: 1, 
		URL: testURL, 
		TargetPrice: 28000},
	}

	repo.On("GetAll", ctx).Return(inputItems, nil)
	scraper.On("FetchCurrentPrice", ctx, inputItems[0].URL).Return(26900.0, nil)
	repo.On("UpdatePrice", ctx, inputItems[0].ID, 26900.0).Return(nil)
	notifier.On("Notify", mock.AnythingOfType("string")).Return(nil)

	// ACT
	svc.CheckAllPrices(ctx)

	// ASSERT (all three methods must call in this case)
	repo.AssertExpectations(t)
	scraper.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

func TestCheckAllPrices_TargetNotReached(t *testing.T) {
	// ARRANGE
	repo := new(MockRepository)
	scraper := new(MockScraper)
	notifier := new(MockNotifier)
	logger := zap.NewNop()

	svc := NewPriceService(repo, scraper, logger, notifier)

	ctx := context.Background()
	inputItems := []domains.Item{
		{ID: 1, URL: testURL, 
		TargetPrice: 20000},
	}

	repo.On("GetAll", ctx).Return(inputItems, nil)
	scraper.On("FetchCurrentPrice", ctx, inputItems[0].URL).Return(26900.0, nil)
	repo.On("UpdatePrice", ctx, inputItems[0].ID, 26900.0).Return(nil)

	// ACT
	svc.CheckAllPrices(ctx)

	// ASSERT (all three methods must call in this case)
	repo.AssertExpectations(t)
	scraper.AssertExpectations(t)
	notifier.AssertNotCalled(t, "Notify")

}

func TestCheckAllPrices_ScraperFailed(t *testing.T) {
	// ARRANGE
	repo := new(MockRepository)
	scraper := new(MockScraper)
	notifier := new(MockNotifier)
	logger := zap.NewNop()

	svc := NewPriceService(repo, scraper, logger, notifier)

	ctx := context.Background()
	inputItems := []domains.Item{
		{ID: 1, URL: testURL, 
		TargetPrice: 28000},
	}

	repo.On("GetAll", ctx).Return(inputItems, nil)
	scraper.On("FetchCurrentPrice", ctx, inputItems[0].URL).Return(0.0, fmt.Errorf("site unavailable"))

	// ACT 
	svc.CheckAllPrices(ctx)

	// ASSERT that last 2 methods of function CheckAllPrices weren't called
	repo.AssertNotCalled(t, "UpdatePrice")
	notifier.AssertNotCalled(t, "Notify")
}