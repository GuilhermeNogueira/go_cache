package go_cache

import (
	"fmt"
	"log"
	"time"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]ItemPriceCache
}

//ItemPriceCache wraps an cached item with the price and it`s expiration.
type ItemPriceCache struct {
	Price      float64
	Expiration int64
}

// return true when expired.
func (item *ItemPriceCache) IsExpired() bool {
	return time.Now().UnixNano() > item.Expiration
}

//planExpiration return when an item will be expired.
func planExpiration(duration time.Duration) int64 {
	return time.Now().Add(duration).UnixNano()
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]ItemPriceCache{},
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	item, ok := c.prices[itemCode]

	if ok {
		if !item.IsExpired() {
			return item.Price, nil
		}
		log.Printf("item [ %v ] expired in cache and will be removed", itemCode)
		delete(c.prices, itemCode) // removing from cache when expired. Not sure if it is necessary.
	}

	log.Printf("fetching item [ %v ] price from external service", itemCode)
	price, err := c.actualPriceService.GetPriceFor(itemCode)

	if err != nil {
		return 0, fmt.Errorf("getting item from service : %v", err.Error())
	}

	c.prices[itemCode] = ItemPriceCache{
		Price:      price,
		Expiration: planExpiration(c.maxAge),
	}

	return price, nil
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {

	var results []float64

	var ch = make(chan float64)
	var errCh = make(chan error)

	defer close(ch)
	defer close(errCh)

	for _, itemCode := range itemCodes {
		go func(code string) {
			price, err := c.GetPriceFor(code)

			if err != nil {
				log.Printf("failed to retrieve price of item [ %v ] for external service", code)
				errCh <- err
			}
			ch <- price
		}(itemCode)
	}

	for {
		select {
		case err := <-errCh:
			log.Printf("operation cancelled due error %v", err)
			return nil, err
		case price := <-ch:
			results = append(results, price)
			if len(results) == len(itemCodes) {
				return results, nil
			}
		}
	}
}
