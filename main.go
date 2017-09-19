package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	PATHS_ENV  = "SSM_PATHS"
	PATHS_FILE = "ssm_paths.txt"
)

var (
	client *ssm.SSM
	paths  []string

	trueBool = true
)

func main() {
	// initialize AWS client
	initClient()

	// initialize path hierarchies
	initPaths()

	// fetch parameters
	params, err := fetchParams(paths)
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

func initPaths() {
	// SSM_PATHS env variable takes precedence
	paths, exists := pathsFromEnv()

	// if SSM_PATHS is not given, read paths from ssm_paths.txt file
	if !exists {
		paths = pathsFromFile(PATHS_FILE)
	}

	fmt.Println(paths)

	// ensure only path hierarchies were given
	for _, path := range paths {
		if !strings.Contains(path, "/") {
			fmt.Printf("Invalid path: '%s' - only path hierarchies are supported (e.g. '/production/webapp/')\n", path)
			os.Exit(1)
		}
	}
}

func pathsFromFile(filename string) []string {
	paths := make([]string, 0)

	f, err := os.Open(filename)
	if err != nil {
		return paths
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}

	return paths
}

func pathsFromEnv() ([]string, bool) {
	var paths []string

	envPaths, found := os.LookupEnv(PATHS_ENV)
	if !found {
		return paths, false
	}

	if envPaths != "" {
		paths = strings.Split(envPaths, ",")
	}
	return paths, true
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
