package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gametimesf/aws-ssm-env/fetch"
)

var (
	client *ssm.SSM
	paths  []string
	tags   []string

	trueBool = true
)

func main() {
	// initialize AWS client
	initClient()

	// initialize command line flags
	initFlags()

	// fetch parameters
	params, err := fetch.FetchParams(paths, tags)
	if err != nil {
		panic(err)
	}

	// print as env variables
	printParams(params)
}

func initClient() {
	session := session.Must(session.NewSession())
	client = ssm.New(session)
}

func initFlags() {
	pathsFlag := flag.String("paths", "", "comma delimited string of parameter path hierarchies")
	tagsFlag := flag.String("tags", "", "comma delimited string of tags to filter by")
	flag.Parse()

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


func getParamNameValues(params []*ssm.Parameter) map[string]string {
	paramVals := make(map[string]string, len(params))
	for _, param := range params {
		split := strings.Split(*param.Name, "/")
		name := split[len(split)-1]
		paramVals[strings.ToUpper(name)] = *param.Value
	}
	return paramVals
}