package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/go-chi/chi"
	"golang.org/x/net/context"
)

type Item struct {
	Month     int       `bigquery:"month"`
	Day       int       `bigquery:"day"`
	Hour      int       `bigquery:"hour"`
	Timestamp time.Time `bigquery:"timestamp"`
	Path      string    `bigquery:"path"`
	Title     string    `bigquery:"title"`
}

func HandleMoshimoshi(w http.ResponseWriter, r *http.Request) {
	sitename := chi.URLParam(r, "site")
	path := chi.URLParam(r, "path")
	title := r.URL.Query().Get("title")

	// create client
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "moshi-moshi-3373")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	defer client.Close()

	// ready table
	t := time.Now()
	tablename := t.Format("2006")
	err = createTableExplicitSchema(ctx, client, sitename, tablename)
	if err != nil {
		// table already created
		log.Printf("%v", err)
	}

	// ready item
	month, _ := strconv.Atoi(t.Format("01"))
	day := t.Day()
	hour := t.Hour()
	replacedPath := strings.Replace(path, "-", "/", -1)
	buf := bytes.NewBufferString("a-know.hateblo.jp/entry/")
	buf.WriteString(replacedPath)
	fullpath := buf.String()
	items := []*Item{
		{Timestamp: t, Path: fullpath, Title: title, Month: month, Day: day, Hour: hour},
	}

	// insert
	u := client.Dataset(sitename).Table(tablename).Uploader()
	err = u.Put(ctx, items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}

	// Pixela increment
	if title != "ExternalMonitoring" {
		req, err := http.NewRequest("PUT", "https://pixe.la/v1/users/a-know-blog/graphs/page-views/increment", nil)
		req.Header.Add("X-USER-TOKEN", os.Getenv("PIXELA_USER_TOKEN"))
		if err != nil {
			log.Printf("Failed to init request. But continue...: %s", err.Error())
		} else {
			client := new(http.Client)
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Failed to send purge request. But continue...: %s", err.Error())
			} else {
				log.Printf("Success to send purge request. continue")
				defer resp.Body.Close()
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "")
}

// func createDataset(ctx context.Context, client *bigquery.Client, datasetID string) error {
// 	meta := &bigquery.DatasetMetadata{
// 		Location: "US", // Create the dataset in the US.
// 	}
// 	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
// 		return err
// 	}
// 	// [END bigquery_create_dataset]
// 	return nil
// }

func createTableExplicitSchema(ctx context.Context, client *bigquery.Client, datasetID, tableID string) error {
	// [START bigquery_create_table]
	sampleSchema := bigquery.Schema{
		{Name: "month", Type: bigquery.IntegerFieldType},
		{Name: "day", Type: bigquery.IntegerFieldType},
		{Name: "hour", Type: bigquery.IntegerFieldType},
		{Name: "title", Type: bigquery.StringFieldType},
		{Name: "path", Type: bigquery.StringFieldType},
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
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
