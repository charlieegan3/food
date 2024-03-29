package main

import (
	"context"
	"log"

	"github.com/charlieegan3/toolbelt/pkg/database"
	"github.com/charlieegan3/toolbelt/pkg/tool"
	"github.com/spf13/viper"

	foodRefreshTool "github.com/charlieegan3/food/pkg/tool"
)

func main() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}

	params := viper.GetStringMapString("database.params")
	connectionString := viper.GetString("database.connectionString")
	db, err := database.Init(connectionString, params, params["dbname"], false)
	if err != nil {
		log.Fatalf("failed to init DB: %s", err)
	}

	tb := tool.NewBelt()
	tb.SetConfig(viper.GetStringMap("tools"))
	tb.SetDatabase(db)

	err = tb.AddTool(&foodRefreshTool.Food{})
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	tb.RunJobs(context.Background())
}
