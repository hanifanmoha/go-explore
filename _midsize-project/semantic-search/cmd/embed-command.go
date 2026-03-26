package cmd

import (
	"fmt"
	"log"

	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/pkg/helper"
	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/services"
	"github.com/spf13/cobra"
)

var EmbedCmd = &cobra.Command{
	Use:   "embed",
	Short: "Generate embeddings for a given text",
	Long:  `Generate embeddings for a given text using the Gemini API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running embeddings...")
		processEmbed()
	},
}

var datasourceFlag string

func init() {
	EmbedCmd.Flags().StringVarP(&datasourceFlag, "datasource", "d", "", "Specify the datasource to embed (e.g., continents, animals). If not specified, embeds all.")
}

func processEmbed() {

	helper.LoadEnv()

	genaiSvc, err := services.NewGenAIService()
	if err != nil {
		log.Println("Error initializing GenAI service:", err)
		return
	}

	dbConnection, err := services.NewDatabaseService()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}
	defer dbConnection.Close()

	datasourceSvc, err := services.NewDataSourceService()
	if err != nil {
		log.Println("Error loading data source:", err)
		return
	}

	if datasourceFlag == "" {
		log.Printf("Please specify a datasource using the --datasource flag")
		return
	}

	topic, err := datasourceSvc.GetDataSource(datasourceFlag)
	if err != nil {
		log.Printf("Error retrieving datasource '%s': %v\n", datasourceFlag, err)
		return
	}

	for _, item := range topic.Items {
		log.Printf("Processing item: %s\n", item.Name)

		vector, err := genaiSvc.GenerateEmbedding(item.Description, services.TaskTypeRetrievalDocument)
		if err != nil {
			log.Printf("Error generating embedding for item %s: %v\n", item.Name, err)
			continue
		}

		err = dbConnection.InsertEmbedding(topic.TopicName, item.Name, item.Description, vector)
		if err != nil {
			log.Printf("Error inserting embedding for item %s into database: %v\n", item.Name, err)
			continue
		}
	}
}
