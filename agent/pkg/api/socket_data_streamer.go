package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/up9inc/mizu/agent/pkg/dependency"
	"github.com/up9inc/mizu/shared/logger"
	tapApi "github.com/up9inc/mizu/tap/api"
	"net/http"
	"time"
)

type EntryStreamer interface {
	Get(ctx context.Context, socketId int, params *WebSocketParams) error
}

type BasenineEntryStreamer struct{}

func (e *BasenineEntryStreamer) Get(ctx context.Context, socketId int, params *WebSocketParams) error {
	entryStreamerSocketConnector := dependency.GetInstance(dependency.EntryStreamerSocketConnector).(EntryStreamerSocketConnector)

	transport := http.DefaultTransport
	tlsClientConfig := &tls.Config{InsecureSkipVerify: true}
	transport.(*http.Transport).TLSClientConfig = tlsClientConfig
	cfg := elasticsearch.Config{
		Addresses: []string{"http://elasticsearch-master.elasticsearch.svc.cluster.local:9200"},
		Transport: transport,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Log.Errorf("Failed to initialize elastic client %v", err)
		return fmt.Errorf("failed to initialize elastic client %v", err)
	}

	res, err := es.Info()
	if err != nil {
		logger.Log.Errorf("Elastic client.Info() ERROR: %v", err)
		return fmt.Errorf("elastic client.Info() ERROR: %v", err)
	}

	defer res.Body.Close()
	if res.IsError() {
		logger.Log.Errorf("Elastic client.Info() ERROR: %v", res.String())
		return fmt.Errorf("elastic client.Info() ERROR: %v", res.String())
	}

	logger.Log.Infof("Starting to stream entries from elastic %v", res.String())

	go func() {
		lastTimestamp := time.Now()

		for {
			var buf bytes.Buffer
			query := map[string]interface{}{
				"query": map[string]interface{}{
					"range": map[string]interface{}{
						"insertionTime": map[string]interface{}{
							"gt": lastTimestamp,
						},
					},
				},
			}
			if err := json.NewEncoder(&buf).Encode(query); err != nil {
				logger.Log.Errorf("Error encoding query: %s", err)
				continue
			}

			// Perform the search request.
			res, err = es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex("mizu_traffic_http"),
				es.Search.WithBody(&buf),
				es.Search.WithTrackTotalHits(true),
				es.Search.WithPretty(),
				es.Search.WithSize(10000),
			)
			if err != nil {
				logger.Log.Errorf("Error getting response: %s", err)
				continue
			}

			if res.IsError() {
				logger.Log.Errorf("Error getting response: %s", res.String())
				continue
			}

			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				logger.Log.Errorf("Error parsing the response body: %s", err)
				continue
			}

			res.Body.Close()

			for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
				source := hit.(map[string]interface{})["_source"]
				jsonValue, _ := json.Marshal(source)

				var entry tapApi.Entry
				if err := json.Unmarshal(jsonValue, &entry); err != nil {
					logger.Log.Errorf("Error parsing the hit: %s", err)
					continue
				}

				entry.Id = hit.(map[string]interface{})["_id"].(string)

				if err := entryStreamerSocketConnector.SendEntry(socketId, &entry, params); err != nil {
					return
				}

				var entryInsertionTime time.Time
				entryInsertionTimeJsonValue, _ := json.Marshal(source.(map[string]interface{})["insertionTime"])
				if err := json.Unmarshal(entryInsertionTimeJsonValue, &entryInsertionTime); err != nil {
					logger.Log.Errorf("Error parsing the hit: %s", err)
					continue
				}

				if entryInsertionTime.After(lastTimestamp) {
					lastTimestamp = entryInsertionTime
				}
			}
		}
	}()

	return nil
}
