### How to use

#### Simulate
`go run cmd/simulate/main.go --csv-file=simulation.csv --simulation-time=24h --spread-percent=0.02 --object-count=1000 --overwrite-csv-file`


#### Plot the Histogram
`go run cmd/graph/main.go --csv-file=simulation.csv --image-file=out.png --graph-start-time=4m --graph-length=4h --overwrite-image-file`

