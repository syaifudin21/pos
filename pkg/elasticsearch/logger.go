package elasticsearch

import (
	"context"
	"os"
	"time"

	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type APILog struct {
	Timestamp  time.Time              `json:"@timestamp"`
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	Status     int                    `json:"status"`
	DurationMs int64                  `json:"duration_ms"`
	User       string                 `json:"user,omitempty"`
	Error      string                 `json:"error,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

var esClient *elasticsearch.Client

func InitElasticLogger() {
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTIC_URL")},
		Username:  os.Getenv("ELASTIC_USER"),
		Password:  os.Getenv("ELASTIC_PASS"),
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Gagal inisialisasi elasticsearch client: %v", err)
	}
	// test connection
	_, err = client.Info()
	if err != nil {
		log.Fatalf("Tidak bisa konek ke elasticsearch: %v", err)
	}
	esClient = client
	log.Println("Elasticsearch logger siap!")
}

func LogAPI(indexName string, logData APILog) {
	if esClient == nil {
		log.Println("Elasticsearch client belum diinisialisasi")
		return
	}
	logData.Timestamp = time.Now().UTC()
	body, err := json.Marshal(logData)
	if err != nil {
		log.Printf("Gagal marshal log: %v", err)
		return
	}
	
	res, err := esClient.Index(
		indexName,
		strings.NewReader(string(body)),
		esClient.Index.WithContext(context.Background()),
		esClient.Index.WithRefresh("false"),
	)
	if err != nil {
		log.Printf("Gagal kirim log ke elastic: %v", err)
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("Elasticsearch response error: %s", res.String())
	}
}
