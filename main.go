package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jamietsao/aws-ssm-env/fetch"
)

var (
	paths []string
	tags  []string
	debug bool
)

func main() {

	// initialize command line flags
	initFlags()

	// fetch parameters
	start := time.Now()
	fetcher := fetch.NewFetcher(os.Getenv("SSM_REGION"), debug)
	params, err := fetcher.FetchParams(paths, tags)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)

	debugf("Params retrieved in %s\n", elapsed)

	// print as env variables
	printParams(params)
}

func initFlags() {
	pathsFlag := flag.String("paths", "", "comma delimited string of parameter path hierarchies (optional)")
	tagsFlag := flag.String("tags", "", "comma delimited string of tags to filter by (optional)")
	flag.BoolVar(&debug, "debug", false, "Enables debug logging when set to true")
	flag.Parse()

	if *pathsFlag == "" && *tagsFlag == "" {
		fmt.Print("Flag required: Either --paths or --tags is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	initPaths(pathsFlag)
	initTags(tagsFlag)
}

func initPaths(pathsFlag *string) {
	if *pathsFlag != "" {
		paths = strings.Split(*pathsFlag, ",")
	}

	// ensure only path hierarchies were given
	for _, path := range paths {
		if !strings.Contains(path, "/") {
			fmt.Printf("Invalid path: '%s' - only path hierarchies are supported (e.g. '/production/webapp/')\n", path)
			os.Exit(1)
		}
	}
}

func initTags(tagsFlag *string) {
	if *tagsFlag != "" {
		tags = strings.Split(*tagsFlag, ",")
	}
}

func printParams(params []*ssm.Parameter) {
	for _, param := range params {
		split := strings.Split(*param.Name, "/")
		name := split[len(split)-1]
		fmt.Printf("%s=%s\n", strings.ToUpper(name), *param.Value)
	}
}

func debugf(format string, a ...interface{}) {
	if debug {
		fmt.Printf("DEBUG -- "+format, a...)
	}
}
