package jobs

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// Refresh will update the status in the database
type Refresh struct {
	DB *sql.DB

	ScheduleOverride string

	DropboxToken string
	DropboxPath  string

	GitHubToken string
	GitHubURL   string
}

func (r *Refresh) Name() string {
	return "refresh"
}

func (r *Refresh) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	goquDB := goqu.New("postgres", r.DB)

	go func() {
		// get the current content_hash
		sel := goquDB.From("food.data").Select("value").Where(goqu.C("key").Eq("content_hash")).Limit(1)
		var currentContentHash string
		_, err := sel.Executor().ScanVal(&currentContentHash)
		if err != nil {
			errCh <- fmt.Errorf("failed to get current status: %w", err)
			return
		}

		config := dropbox.Config{
			Token: r.DropboxToken,
		}
		dbx := files.New(config)
		metadata, err := dbx.GetMetadata(&files.GetMetadataArg{
			Path: r.DropboxPath,
		})
		if err != nil {
			errCh <- fmt.Errorf("failed to get metadata: %w", err)
			return
		}

		fileMetadata, ok := metadata.(*files.FileMetadata)
		if !ok {
			errCh <- fmt.Errorf("metadata was an unexpected type: %T", metadata)
			return
		}

		if currentContentHash != fileMetadata.ContentHash {
			log.Println("updated needed")

			jobData := map[string]interface{}{
				"event_type":     "refresh",
				"client_payload": map[string]interface{}{},
			}

			bodyJSON, err := json.Marshal(jobData)
			if err != nil {
				errCh <- fmt.Errorf("failed to marshal json: %w", err)
				return
			}

			client := http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("POST", r.GitHubURL, bytes.NewReader(bodyJSON))
			if err != nil {
				errCh <- fmt.Errorf("failed to create request: %w", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+r.GitHubToken)

			resp, err := client.Do(req)
			if err != nil {
				errCh <- fmt.Errorf("failed to make request: %w", err)
				return
			}

			if resp.StatusCode >= 400 {
				errCh <- fmt.Errorf("failed to make request: %d", resp.StatusCode)
				return
			}

			// update the content_hash in the database
			ins := goquDB.Insert("food.data").Rows(goqu.Record{"key": "content_hash", "value": fileMetadata.ContentHash}).
				OnConflict(goqu.DoUpdate("key", goqu.Record{"value": fileMetadata.ContentHash}))
			_, err = ins.Executor().Exec()
			if err != nil {
				errCh <- fmt.Errorf("failed to update content hash record: %w", err)
				return
			}
		}

		doneCh <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errCh:
		return fmt.Errorf("job failed with error: %s", e)
	case <-doneCh:
		return nil
	}
}

func (r *Refresh) Timeout() time.Duration {
	return 15 * time.Second
}

func (r *Refresh) Schedule() string {
	if r.ScheduleOverride != "" {
		return r.ScheduleOverride
	}
	return "0 */5 * * * *"
}
