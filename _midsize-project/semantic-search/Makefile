.PHONY: run-embeddings run-search run-serve

# make run-embeddings DATASOURCE=datasource_name
run-embeddings:
	go run main.go embed $(if $(DATASOURCE),--datasource $(DATASOURCE))

run-search:
	go run main.go search

run-serve:
	go run main.go serve