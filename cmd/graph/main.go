package main

import (
	"fmt"
	"os"

	"github.com/Tomasz-Smelcerz-SAP/jitter/cmd"
	"github.com/Tomasz-Smelcerz-SAP/jitter/internal/draw"
	"github.com/Tomasz-Smelcerz-SAP/jitter/internal/histogram"
	"github.com/Tomasz-Smelcerz-SAP/jitter/internal/model"
)

const (
	defaultArgGraphStartTime = "24h"
	defaultArgGraphLength    = "60m"
	defaultArgImageFileName  = "out.png"
)

func main() {

	options := parseCLIArguments(os.Args)

	fileAlreadyExists, err := cmd.FileExists(options.imageFileName)
	if err != nil {
		fmt.Println("Error checking if image file exists:", err)
		os.Exit(1)
	}
	if fileAlreadyExists {
		if !options.overwriteImageFile {
			fmt.Printf("Image file already exists: %s\n", options.imageFileName)
			os.Exit(1)
		}
	}

	fmt.Println("================================================================================")
	fmt.Println("Reding input data from CSV file...")
	file, err := os.Open(options.csvFileName)
	objects, err := model.UnmarshalObjSet(file)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	objCount := len(objects)
	fmt.Println("   Read", objCount, "objects")

	fmt.Println("================================================================================")
	fmt.Println("Calculating the histogram...")
	var graphStartTimeMillis int = cmd.SecondsToMillis(options.graphStartTimeSeconds)
	var graphLengthMillis int = cmd.SecondsToMillis(options.graphLengthSeconds)
	bucketCount := 1000

	if graphLengthMillis%bucketCount != 0 {
		panic("graphLengthMillis must be divisible by bucketCount")
	}

	timePerBucket := graphLengthMillis / bucketCount // In this case division is always possible!
	hist := histogram.NewHistogram(graphStartTimeMillis, timePerBucket, bucketCount)
	for i := 0; i < objCount; i++ {
		obj := objects[i]
		for _, schedule := range obj.Schedules() {
			if schedule >= float64(graphStartTimeMillis) && schedule < float64(graphStartTimeMillis+graphLengthMillis) {
				hist.AddDataPoint(int(schedule))
			}
		}
	}

	expectedSchedules := float64(objCount) / model.AverageScheduleTime * float64(graphLengthMillis) // Assuming perfectly uniform distribution
	fmt.Println("   Expected schedules:", int(expectedSchedules))
	fmt.Println("   Total schedules:", hist.TotalCount())

	fmt.Println("================================================================================")
	fmt.Println("Drawing histogram")
	draw.Draw(hist, options.argGraphStartTime, options.argGraphLength, options.imageFileName)

	fmt.Println("================================================================================")
	fmt.Println("Done")
}

func parseCLIArguments(osArgs []string) options {
	res := options{}

	if len(osArgs) < 2 {
		fmt.Println("Reads the simulation data file and plots results as a histogram with configurable time window.")
		fmt.Println("Usage: go run . --csv-file=<path> [--image-file=<path>] [--overwrite-image-file] --graph-start-time=<time> --graph-length=<time>")
		fmt.Println("Example: go run . --csv-file=simulation.csv --image-file=out.png --graph-start-time=4m --graph-length=4h")
		os.Exit(1)
	}

	args := cmd.Arguments{}
	for i := 1; i < len(osArgs); i++ {
		args.Add(osArgs[i])
	}

	argCSVFileName, ok := args.Get("--csv-file")
	if !ok {
		fmt.Println("Missing argument --csv-file")
		os.Exit(1)
	}
	res.csvFileName = argCSVFileName

	argImageFileName, ok := args.Get("--image-file")
	if !ok {
		argImageFileName = defaultArgImageFileName
	}
	res.imageFileName = argImageFileName

	_, ok = args.Get("--overwrite-image-file")
	res.overwriteImageFile = ok

	argGraphStartTime, ok := args.Get("--graph-start-time")
	if !ok {
		argGraphStartTime = defaultArgGraphStartTime
	}
	res.argGraphStartTime = argGraphStartTime
	graphStartTimeSeconds, err := cmd.AsSeconds(argGraphStartTime)
	if err != nil {
		fmt.Printf("Invalid argument value for --graph-start-time: %s\n", argGraphStartTime)
		os.Exit(1)
	}
	res.graphStartTimeSeconds = graphStartTimeSeconds

	argGraphLength, ok := args.Get("--graph-length")
	if !ok {
		argGraphLength = defaultArgGraphLength
	}
	res.argGraphLength = argGraphLength
	graphLengthSeconds, err := cmd.AsSeconds(argGraphLength)
	if err != nil {
		fmt.Printf("Invalid argument value for --graph-length: %s\n", argGraphLength)
		os.Exit(1)
	}
	res.graphLengthSeconds = graphLengthSeconds

	return res
}

type options struct {
	csvFileName           string
	imageFileName         string
	overwriteImageFile    bool
	argGraphStartTime     string
	graphStartTimeSeconds int
	argGraphLength        string
	graphLengthSeconds    int
}
