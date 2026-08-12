package main

import (
	"fmt"
	"time"
)

type Cache interface {
	Get(key string) (string, bool) // bool = found
	Set(key string, value string, ttl time.Duration)
	Delete(key string)
	Clear()
}

type entry struct {
	value     string
	expiresAt time.Time
}

type InMemoryCache struct {
	data map[string]entry
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{data: map[string]entry{}}
}

func (c *InMemoryCache) Get(key string) (string, bool) {
	item, ok := c.data[key]
	if !ok {
		return "", false
	}
	if time.Now().After(item.expiresAt) {
		delete(c.data, key)
		return "", false
	}
	return item.value, true
}

func (c *InMemoryCache) Set(key string, value string, ttl time.Duration) {
	item := entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.data[key] = item
}

func (c *InMemoryCache) Delete(key string) {
	delete(c.data, key)
}

func (c *InMemoryCache) Clear() {
	c.data = make(map[string]entry)
}

type NoCache struct{}

func (n *NoCache) Get(key string) (string, bool)                   { return "", false }
func (n *NoCache) Set(key string, value string, ttl time.Duration) {}
func (n *NoCache) Delete(key string)                               {}
func (n *NoCache) Clear()                                          {}

type WeatherService struct {
	cache     Cache
	callCount int
}

func NewWeatherService(cache Cache) *WeatherService {
	return &WeatherService{cache: cache, callCount: 0}
}

func (w *WeatherService) GetWeather(city string) string {
	key := "weather:" + city
	cache, found := w.cache.Get(key)
	if found {
		return fmt.Sprintf("%s: %s (served from cache)", city, cache)
	}
	weatherData := map[string]string{
		"Karachi": "Sunny, 35°C",
		"Lahore":  "Cloudy, 28°C",
	}
	raw, ok := weatherData[city]
	if !ok {
		raw = "Unknown"
	}

	w.callCount++
	w.cache.Set(key, raw, 5*time.Minute)
	return fmt.Sprintf("%s: %s (fetched from API)", city, raw)
}

func (w *WeatherService) Stats() string {
	return fmt.Sprintf("API calls made: %d", w.callCount)
}

func main() {
	fmt.Println("With cache:")
	cache := NewInMemoryCache()
	service := NewWeatherService(cache)
	fmt.Println(service.GetWeather("Karachi"))
	fmt.Println(service.GetWeather("Lahore"))
	fmt.Println(service.GetWeather("Karachi"))
	fmt.Println(service.GetWeather("Karachi"))
	fmt.Println(service.Stats())
	fmt.Println("\nWith NoCache:")
	noCacheService := NewWeatherService(&NoCache{})
	fmt.Println(noCacheService.GetWeather("Karachi"))
	fmt.Println(noCacheService.GetWeather("Karachi"))
	fmt.Println(noCacheService.GetWeather("Karachi"))
	fmt.Println(noCacheService.Stats())
}
