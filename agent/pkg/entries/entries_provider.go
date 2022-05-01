package entries

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"net/http"
	"time"

	basenine "github.com/up9inc/basenine/client/go"
	"github.com/up9inc/mizu/agent/pkg/app"
	"github.com/up9inc/mizu/agent/pkg/har"
	"github.com/up9inc/mizu/agent/pkg/models"
	"github.com/up9inc/mizu/shared"
	"github.com/up9inc/mizu/shared/logger"
	tapApi "github.com/up9inc/mizu/tap/api"
)

type EntriesProvider interface {
	GetEntries(entriesRequest *models.EntriesRequest) ([]*tapApi.EntryWrapper, *basenine.Metadata, error)
	GetEntry(singleEntryRequest *models.SingleEntryRequest, entryId string) (*tapApi.EntryWrapper, error)
}

type BasenineEntriesProvider struct{}

func (e *BasenineEntriesProvider) GetEntries(entriesRequest *models.EntriesRequest) ([]*tapApi.EntryWrapper, *basenine.Metadata, error) {
	data, meta, err := basenine.Fetch(shared.BasenineHost, shared.BaseninePort,
		entriesRequest.LeftOff, entriesRequest.Direction, entriesRequest.Query,
		entriesRequest.Limit, time.Duration(entriesRequest.TimeoutMs)*time.Millisecond)
	if err != nil {
		return nil, nil, err
	}

	var dataSlice []*tapApi.EntryWrapper

	for _, row := range data {
		var entry *tapApi.Entry
		err = json.Unmarshal(row, &entry)
		if err != nil {
			return nil, nil, err
		}

		extension := app.ExtensionsMap[entry.Protocol.Name]
		base := extension.Dissector.Summarize(entry)

		dataSlice = append(dataSlice, &tapApi.EntryWrapper{
			Protocol: entry.Protocol,
			Data:     entry,
			Base:     base,
		})
	}

	var metadata *basenine.Metadata
	err = json.Unmarshal(meta, &metadata)
	if err != nil {
		logger.Log.Debugf("Error recieving metadata: %v", err.Error())
	}

	return dataSlice, metadata, nil
}

func (e *BasenineEntriesProvider) GetEntry(singleEntryRequest *models.SingleEntryRequest, entryId string) (*tapApi.EntryWrapper, error) {
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
		return nil, fmt.Errorf("failed to initialize elastic client %v", err)
	}

	res, err := es.Info()
	if err != nil {
		logger.Log.Errorf("Elastic client.Info() ERROR: %v", err)
		return nil, fmt.Errorf("elastic client.Info() ERROR: %v", err)
	}

	defer res.Body.Close()
	if res.IsError() {
		logger.Log.Errorf("Elastic client.Info() ERROR: %v", res.String())
		return nil, fmt.Errorf("elastic client.Info() ERROR: %v", res.String())
	}

	resp, err := es.GetSource("mizu_traffic_http", entryId)
	if err != nil {
		logger.Log.Errorf("Error getting response: %s", err)
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		logger.Log.Errorf("Error getting response: %s", resp.String())
		return nil, fmt.Errorf("error getting response: %s", resp.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		logger.Log.Errorf("Error parsing the response body: %s", err)
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	jsonValue, _ := json.Marshal(r)

	var entry *tapApi.Entry
	if err := json.Unmarshal(jsonValue, &entry); err != nil {
		logger.Log.Errorf("Error parsing the hit: %s", err)
		return nil, fmt.Errorf("error parsing the hit: %s", err)
	}

	extension := app.ExtensionsMap[entry.Protocol.Name]
	base := extension.Dissector.Summarize(entry)
	var representation []byte
	representation, err = extension.Dissector.Represent(entry.Request, entry.Response)
	if err != nil {
		return nil, err
	}

	var rules []map[string]interface{}
	var isRulesEnabled bool
	if entry.Protocol.Name == "http" {
		harEntry, _ := har.NewEntry(entry.Request, entry.Response, entry.StartTime, entry.ElapsedTime)
		_, rulesMatched, _isRulesEnabled := models.RunValidationRulesState(*harEntry, entry.Destination.Name)
		isRulesEnabled = _isRulesEnabled
		inrec, _ := json.Marshal(rulesMatched)
		if err := json.Unmarshal(inrec, &rules); err != nil {
			logger.Log.Error(err)
		}
	}

	return &tapApi.EntryWrapper{
		Protocol:       entry.Protocol,
		Representation: string(representation),
		Data:           entry,
		Base:           base,
		Rules:          rules,
		IsRulesEnabled: isRulesEnabled,
	}, nil
}
