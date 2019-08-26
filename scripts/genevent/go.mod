module genevent

go 1.12

replace github.com/apognu/gocal => ./gocal

replace github.com/apognu/gocal/parser => ./gocal/parser

require (
	github.com/apognu/gocal v0.4.1
	github.com/stretchr/testify v1.4.0 // indirect
)
