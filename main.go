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
	tags   []string

	trueBool = true
)

func main() {
	// initialize AWS client
	initClient()

	// initialize command line flags
	initFlags()

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

func fetchParams(paths []string) ([]*ssm.Parameter, error) {
	// create tag filters
	tagFilters := make([]*ssm.ParameterStringFilter, len(tags))
	for i, tag := range tags {
		tagFilter := fmt.Sprintf("tag:%s", tag)
		tagFilters[i] = &ssm.ParameterStringFilter{
			Key: &tagFilter,
		}
	}

	// TEMP: until parameter-filters work for get-parameters-by-path (https://github.com/aws/aws-cli/issues/2850),
	// 1) retrieve parameters by tags via describe-parameters
	// 2) retrieve parameters by path via get-parameters-by-path
	// 3) calculate union of two sets

	// retrieve all parameters with given tags
	paramNames, err := describeParams(tagFilters)
	if err != nil {
		return nil, err
	}

	// retrieve params for given paths
	params, err := getParamsByPath(paths)
	if err != nil {
		return params, err
	}

	// calculate union of two sets
	union := calcUnion(paramNames, params)

	return union, nil
}

func describeParams(filters []*ssm.ParameterStringFilter) ([]string, error) {
	paramNames := make([]string, 0)

	done := false
	var nextToken string
	for !done {
		input := &ssm.DescribeParametersInput{
			ParameterFilters: filters,
		}

		if nextToken != "" {
			input.SetNextToken(nextToken)
		}

		output, err := client.DescribeParameters(input)
		if err != nil {
			return paramNames, err
		}

		for _, param := range output.Parameters {
			paramNames = append(paramNames, *param.Name)
		}

		// there are more parameters if nextToken is given in response
		if output.NextToken != nil {
			nextToken = *output.NextToken
		} else {
			done = true
		}
	}

	return paramNames, nil
}

func getParamsByPath(paths []string) ([]*ssm.Parameter, error) {
	params := make([]*ssm.Parameter, 0)

	// retrieve params for all paths
	for _, path := range paths {

		done := false
		var nextToken string
		for !done {
			input := &ssm.GetParametersByPathInput{
				Path:           &path,
				Recursive:      &trueBool,
				WithDecryption: &trueBool,
			}

			if nextToken != "" {
				input.SetNextToken(nextToken)
			}

			output, err := client.GetParametersByPath(input)
			if err != nil {
				return params, err
			}

			params = append(params, output.Parameters...)

			// there are more parameters for this path if nextToken is given in response
			if output.NextToken != nil {
				nextToken = *output.NextToken
			} else {
				done = true
			}
		}
	}

	return params, nil
}

func calcUnion(paramNames []string, params []*ssm.Parameter) []*ssm.Parameter {
	// build map lookup
	lookup := make(map[string]bool)
	for _, paramName := range paramNames {
		lookup[paramName] = true
	}

	// calculate union
	union := make([]*ssm.Parameter, 0)
	for _, param := range params {
		if lookup[*param.Name] {
			union = append(union, param)
		}
	}

	return union
}

func printParams(params []*ssm.Parameter) {
	for _, param := range params {
		split := strings.Split(*param.Name, "/")
		name := split[len(split)-1]
		fmt.Printf("%s=%s\n", strings.ToUpper(name), *param.Value)
	}
}
