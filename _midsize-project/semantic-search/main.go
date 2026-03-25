package main

import (
	"fmt"
	"os"

	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "semantic-search",
	Short: "A semantic search application using embeddings",
	Long:  `A CLI tool for semantic search with Gemini API embeddings`,
}

func init() {
	rootCmd.AddCommand(cmd.EmbedCmd)
	rootCmd.AddCommand(cmd.ServeCmd)
	rootCmd.AddCommand(cmd.SearchCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
