#!/bin/bash

go test ./... -coverprofile=coverage.full.out
cat coverage.full.out | grep -v "test/" > coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
go-cover-treemap -coverprofile coverage.out > coverage.svg
