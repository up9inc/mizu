package elastic

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/up9inc/mizu/shared"
	"github.com/up9inc/mizu/shared/logger"
	"github.com/up9inc/mizu/tap/api"
	"net/http"
	"sync"
	"time"
)

type client struct {
	es            *elasticsearch.Client
	index         string
	insertedCount int
}

var instance *client
var once sync.Once

func GetInstance() *client {
	once.Do(func() {
		instance = newClient()
	})
	return instance
}

func (client *client) Configure(config shared.ElasticConfig) {
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
	}

	// Have the client instance return a response
	res, err := es.Info()
	if err != nil {
		logger.Log.Errorf("Elastic client.Info() ERROR: %v", err)
	} else {
		defer res.Body.Close()
		client.es = es
		client.index = "mizu_traffic_http"
		client.insertedCount = 0
		logger.Log.Infof("Elastic client configured, index: %s, cluster info: %v", client.index, res)
	}
}

func newClient() *client {
	return &client{
		es:    nil,
		index: "",
	}
}

func (client *client) PushEntry(entry *api.Entry) {
	if client.es == nil {
		return
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		logger.Log.Errorf("json.Marshal ERROR: %v", err)
		return
	}

	var entryMap map[string]interface{}
	if err := json.Unmarshal(entryJson, &entryMap); err != nil {
		logger.Log.Errorf("json.Unmarshal ERROR: %v", err)
		return
	}

	entryMap["insertionTime"] = time.Now()

	entryJson, err = json.Marshal(entryMap)
	if err != nil {
		logger.Log.Errorf("json.Marshal ERROR: %v", err)
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString(string(entryJson))
	res, _ := client.es.Index(client.index, &buffer)
	if res.StatusCode == 201 {
		client.insertedCount += 1
	}
}
