package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	client *ssm.SSM
	paths  []string

	trueBool = true
)

func main() {
	// parse command line flags
	initFlags()

	// initialize AWS client
	initClient()

	// fetch parameters
	params, err := fetchParams(paths)
	if err != nil {
		panic(err)
	}

	// print as env variables
	printParams(params)
}

func initFlags() {
	pathsFlag := flag.String("paths", "", "comma delimited string of parameter path hierarchies")
	flag.Parse()

	if *pathsFlag != "" {
		paths = strings.Split(*pathsFlag, ",")
	}

	if len(paths) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, path := range paths {
		if !strings.Contains(path, "/") {
			fmt.Printf("Invalid path: '%s' - only path hierarchies are supported (e.g. '/production/eventsdiscovery/')\n", path)
			os.Exit(1)
		}
	}
}

func initClient() {
	session := session.Must(session.NewSession())
	client = ssm.New(session)
}

func fetchParams(paths []string) ([]*ssm.Parameter, error) {
	params := make([]*ssm.Parameter, 0)

	for _, path := range paths {
		resp, err := client.GetParametersByPath(&ssm.GetParametersByPathInput{
			Path:           &path,
			WithDecryption: &trueBool,
		})

		if err != nil {
			return params, err
		}

		params = append(params, resp.Parameters...)
	}

	return params, nil
}

func printParams(params []*ssm.Parameter) {
	for _, param := range params {
		split := strings.Split(*param.Name, "/")
		name := split[len(split)-1]
		fmt.Printf("%s=%s\n", strings.ToUpper(name), *param.Value)
	}
}
