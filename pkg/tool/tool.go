package tool

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/charlieegan3/toolbelt/pkg/apis"
	"github.com/gorilla/mux"

	"github.com/charlieegan3/food/pkg/tool/jobs"
)

//go:embed migrations
var foodDatabaseMigrations embed.FS

// Food is a tool for functions relating to the food website
type Food struct {
	config *gabs.Container
	db     *sql.DB
}

func (f *Food) Name() string {
	return "food"
}

func (f *Food) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Config:   true,
		Jobs:     true,
		Database: true,
	}
}

func (f *Food) SetConfig(config map[string]any) error {
	f.config = gabs.Wrap(config)

	return nil
}
func (f *Food) Jobs() ([]apis.Job, error) {
	var j []apis.Job
	var path string
	var ok bool

	// load all config
	path = "dropbox.token"
	dropboxToken, ok := f.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "dropbox.path"
	dropboxPath, ok := f.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "github.token"
	githubToken, ok := f.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "github.url"
	githubURL, ok := f.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "jobs.refresh.schedule"
	schedule, ok := f.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	return []apis.Job{
		&jobs.Refresh{
			DB:               f.db,
			ScheduleOverride: schedule,
			DropboxToken:     dropboxToken,
			DropboxPath:      dropboxPath,
			GitHubToken:      githubToken,
			GitHubURL:        githubURL,
		},
	}, nil
}
func (f *Food) ExternalJobsFuncSet(fun func(job apis.ExternalJob) error) {}

func (f *Food) DatabaseMigrations() (*embed.FS, string, error) {
	return &foodDatabaseMigrations, "migrations", nil
}
func (f *Food) DatabaseSet(db *sql.DB) {
	f.db = db
}

func (f *Food) HTTPPath() string                    { return "" }
func (f *Food) HTTPHost() string                    { return "" }
func (f *Food) HTTPAttach(router *mux.Router) error { return nil }
