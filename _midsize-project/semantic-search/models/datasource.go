package models

type DataSourceTopic struct {
	TopicName string           `json:"topic"`
	Items     []DataSourceItem `json:"items"`
}

type DataSourceItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
