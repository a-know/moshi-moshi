package handlers

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/go-chi/chi"
	"golang.org/x/net/context"
)

type Item struct {
	Timestamp time.Time `bigquery:"timestamp"`
	Path      string    `bigquery:"path"`
	Title     string    `bigquery:"title"`
}

func HandleMoshimoshi(w http.ResponseWriter, r *http.Request) {
	// sitename := chi.URLParam(r, "site")
	path := chi.URLParam(r, "path")
	title := r.URL.Query().Get("title")

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "moshi-moshi-3373")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	defer client.Close()

	err = createTableExplicitSchema(ctx, client, "testdataset", "table")
	if err != nil {
		// table already created
	}

	u := client.Dataset("testdataset").Table("table").Uploader()
	now := time.Now()
	items := []*Item{
		{Timestamp: now, Path: path, Title: title},
	}
	err = u.Put(ctx, items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "")
}

func createDataset(ctx context.Context, client *bigquery.Client, datasetID string) error {
	meta := &bigquery.DatasetMetadata{
		Location: "US", // Create the dataset in the US.
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		return err
	}
	// [END bigquery_create_dataset]
	return nil
}

func createTableExplicitSchema(ctx context.Context, client *bigquery.Client, datasetID, tableID string) error {
	// [START bigquery_create_table]
	sampleSchema := bigquery.Schema{
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
		{Name: "path", Type: bigquery.StringFieldType},
		{Name: "title", Type: bigquery.StringFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema: sampleSchema,
		// ExpirationTime: time.Now().AddDate(1, 0, 0), // Table will be automatically deleted in 1 year.
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	// [END bigquery_create_table]
	return nil
}
