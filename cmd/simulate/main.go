package main

import (
	"fmt"

	"github.com/Tomasz-Smelcerz-SAP/jitter/cmd"
	"github.com/Tomasz-Smelcerz-SAP/jitter/internal/model"

	"math/rand/v2"
	"os"
	"strconv"
)

const (
	defaultArgSimulationTime = "24h"
	defaultSpreadPercent     = "0.02" // 2%
	defaultObjectCount       = "1000"
)

func main() {

	opts := parseCLIArguments(os.Args)

	if opts.objCount <= 100 || opts.objCount > 100000 {
		fmt.Printf("Object count must be between 100 and 100000\n")
		os.Exit(1)
	}

	fileExists, err := cmd.FileExists(opts.csvFileName)
	if err != nil {
		fmt.Println("Error checking if csv file exists:", err)
		os.Exit(1)
	}
	if fileExists {
		if !opts.overwriteCsvFile {
			fmt.Printf("File %s already exists. Please remove it or choose another file name.\n", opts.csvFileName)
			os.Exit(1)
		}
	}

	fmt.Println("================================================================================")
	fmt.Println("Generating the scheduling of objects over time:")
	fmt.Printf("   Simulation time: %d:%d:%d [h:m:s]\n", opts.simulationTimeSeconds/3600, (opts.simulationTimeSeconds%3600)/60, opts.simulationTimeSeconds%60)
	fmt.Printf("   Object count: %d\n", opts.objCount)
	fmt.Printf("   Spread percent: %.2f\n", opts.spreadPercent)

	var simulationTimeMillis int = cmd.SecondsToMillis(opts.simulationTimeSeconds)
	var initialScheduleMillis int = cmd.MinutesToMillis(5)

	var objects = model.ObjSet{}

	fmt.Println("================================================================================")
	fmt.Println("Initializing objects...")
	rs := model.RandomSupport{
		Float64: rand.Float64,
	}

	for i := 0; i < opts.objCount; i++ {
		obj := model.NewObject(i, float64(initialScheduleMillis), opts.spreadPercent).SetRandomSupport(rs)
		objects = append(objects, obj)
	}

	fmt.Println("================================================================================")
	fmt.Println("Simulating re-schedules...")
	for i := 0; i < opts.objCount; i++ {
		obj := objects[i]
		// Simulate re-schedules
		for obj.LastSchedule() < float64(simulationTimeMillis) {
			obj.AddRandomSchedule()
		}
	}

	fmt.Println("================================================================================")
	fmt.Println("Writing object schedules to a file...")

	file, err := os.Create(opts.csvFileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing CSV file: %v\n", err)
		}
	}()

	objects.Marshal(file)

	fmt.Println("================================================================================")
	fmt.Println("Done")
}

func parseCLIArguments(osArgs []string) options {

	res := options{}
	if len(osArgs) < 2 {
		fmt.Println("Runs the simulation and stores the results in a CSV file.")
		fmt.Println("Usage: go run . --csv-file=<path> [--simulation-time=<time>] [--spread-percent=<float>] [--object-count=<uint>] [--overwrite-csv-file]")
		fmt.Println("Example: go run . --csv-file=simulation.csv --simulation-time=24h --spread-percent=0.02 --object-count=1000")
		os.Exit(1)
	}

	args := cmd.Arguments{}
	for i := 1; i < len(os.Args); i++ {
		args.Add(os.Args[i])
	}

	csvFileName, ok := args.Get("--csv-file")
	if !ok {
		fmt.Println("Missing argument --csv-file")
		os.Exit(1)
	}
	res.csvFileName = csvFileName

	_, ok = args.Get("--overwrite-csv-file")
	res.overwriteCsvFile = ok

	argSimulationTime, ok := args.Get("--simulation-time")
	if !ok {
		argSimulationTime = defaultArgSimulationTime
	}
	simulationTimeSeconds, err := cmd.AsSeconds(argSimulationTime)
	if err != nil {
		fmt.Printf("Invalid argument value for --simulation-time: %s\n", argSimulationTime)
		os.Exit(1)
	}
	res.simulationTimeSeconds = simulationTimeSeconds

	argSpreadPercent, ok := args.Get("--spread-percent")
	if !ok {
		argSpreadPercent = defaultSpreadPercent
	}
	spreadPercent, err := strconv.ParseFloat(argSpreadPercent, 64)
	if err != nil {
		fmt.Printf("Invalid argument value for --spread-percent: %s\n", argSpreadPercent)
		os.Exit(1)
	}
	res.spreadPercent = spreadPercent

	argObjectCount, ok := args.Get("--object-count")
	if !ok {
		argObjectCount = defaultObjectCount
	}
	objCount, err := strconv.Atoi(argObjectCount)
	if err != nil {
		fmt.Printf("Invalid argument value for --object-count: %s\n", argObjectCount)
		os.Exit(1)
	}
	res.objCount = objCount

	return res
}

type options struct {
	csvFileName           string
	overwriteCsvFile      bool
	simulationTimeSeconds int
	spreadPercent         float64
	objCount              int
}
