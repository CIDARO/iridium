package metrics

import (
	"encoding/json"
	"net/http"

	"github.com/CIDARO/iridium/internal/db"
	"github.com/CIDARO/iridium/internal/utils"
	"github.com/bradfitz/gomemcache/memcache"
)

// Metric struct
type Metric struct {
	GetCount         int64 `json:"get_count"`
	PostCount        int64 `json:"post_count"`
	PutCount         int64 `json:"put_count"`
	DeleteCount      int64 `json:"delete_count"`
	OtherCount       int64 `json:"other_count"`
	Count            int64 `json:"count"`
	AvgContentLength int64 `json:"avg_content_length"`
	MaxContentLength int64 `json:"max_content_length"`
	MinContentLength int64 `json:"min_content_length"`
}

// StoreMetrics for the given request
func StoreMetrics(req *http.Request) {
	// default values
	getCount := int64(0)
	postCount := int64(0)
	putCount := int64(0)
	deleteCount := int64(0)
	otherCount := int64(0)
	count := int64(1)
	avgContentLength := req.ContentLength
	maxContentLength := req.ContentLength
	minContentLength := req.ContentLength

	item, err := db.Memcache.Get(req.URL.String())

	if err != nil {
		var metrics Metric

		err := json.Unmarshal(item.Value, &metrics)
		if err != nil {
			utils.Logger.Errorf("error while unmarshaling metrics: %v", err)
			return
		}

		getCount = metrics.GetCount
		postCount = metrics.PostCount
		putCount = metrics.PostCount
		deleteCount = metrics.DeleteCount
		otherCount = metrics.OtherCount
		count = metrics.Count
		avgContentLength = (metrics.AvgContentLength + req.ContentLength) / 2
		maxContentLength = metrics.MaxContentLength
		minContentLength = metrics.MinContentLength

		if req.ContentLength > maxContentLength {
			maxContentLength = req.ContentLength
		} else if req.ContentLength < minContentLength {
			minContentLength = req.ContentLength
		}
	}

	switch req.Method {
	case "GET":
		getCount++
	case "POST":
		postCount++
	case "PUT":
		putCount++
	case "DELETE":
		deleteCount++
	default:
		otherCount++
	}

	parsed, err := json.Marshal(Metric{
		GetCount:         getCount,
		PostCount:        postCount,
		PutCount:         putCount,
		DeleteCount:      deleteCount,
		Count:            count,
		AvgContentLength: avgContentLength,
		MaxContentLength: maxContentLength,
		MinContentLength: minContentLength,
	})

	if err != nil {
		utils.Logger.Errorf("error while marshaling metrics: %v", err)
		return
	}

	db.Memcache.Set(&memcache.Item{
		Key:   req.URL.String(),
		Value: parsed,
	})
}
