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

	for _, topic := range datasourceSvc.Topics {
		log.Printf("Process embedding for topic: %s with %d items\n", topic.TopicName, len(topic.Items))

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

}
