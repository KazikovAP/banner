package cache

import (
	"banner/internal/models"
	"strconv"
	"sync"
	"time"
)

type Cache struct {
	Banner    models.Banner
	UpdatedAt time.Time
}

type CacheBanner struct {
	Banners map[string]Cache
	sync.RWMutex
}

var cache CacheBanner

func init() {
	cache = CacheBanner{
		Banners: make(map[string]Cache),
	}
}

func GenerateCacheKey(featureID, tagID int) string {
	return strconv.Itoa(featureID) + "-" + strconv.Itoa(tagID)
}

func GetBannerFromCache(featureID, tagID int) (*models.Banner, bool) {
	cache.RLock()
	defer cache.RUnlock()

	key := GenerateCacheKey(featureID, tagID)
	cached, found := cache.Banners[key]
	if !found {
		return nil, false
	}

	if time.Since(cached.UpdatedAt) > 5*time.Minute {
		delete(cache.Banners, key)
		return nil, false
	}

	return &cached.Banner, true
}

func StorageBannerInCache(featureID, tagID int, banner models.Banner) {
	cache.Lock()
	defer cache.Unlock()

	key := GenerateCacheKey(featureID, tagID)
	cache.Banners[key] = Cache{
		Banner:    banner,
		UpdatedAt: time.Now(),
	}
}
