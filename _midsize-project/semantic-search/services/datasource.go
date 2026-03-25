package services

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/models"
)

var DataSourceFilePath = []string{
	"./datasource/continents.json",
	"./datasource/planets.json",
}

type DataSourceService struct {
	Topics []models.DataSourceTopic
}

func NewDataSourceService() (*DataSourceService, error) {

	ds := &DataSourceService{}
	err := ds.LoadDataSource()

	return ds, err
}

func (d *DataSourceService) LoadDataSource() error {

	for _, path := range DataSourceFilePath {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		var topic models.DataSourceTopic
		err = decoder.Decode(&topic)
		if err != nil {
			return err
		}

		d.Topics = append(d.Topics, topic)
		log.Printf("Loaded topic: %s with %d items\n", topic.TopicName, len(topic.Items))
	}

	return nil
}
