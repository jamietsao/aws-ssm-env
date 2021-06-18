package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

const (
	defaultBackoff time.Duration = 10 * time.Second

	envBackoff = "AWS_SSM_ENV_BACKOFF"
)

var (
	client *ssm.Client
	paths  []string
	tags   []string
)

func main() {
	// initialize AWS client
	err := initClient()
	if err != nil {
		log.Fatalf("error initializing client: %v", err)
		os.Exit(1)
	}

	// initialize command line flags
	initFlags()

	// fetch parameters
	params, err := fetchParams(paths)
	if err != nil {
		log.Fatalf("error fetching parameters: %v", err)
		os.Exit(2)
	}

	// print as env variables
	printParams(params)
}

func initClient() error {
	// Create base config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	// Create retryer with custom backoff if env var present
	customRetryer := retry.NewStandard(func(o *retry.StandardOptions) {
		o.MaxBackoff = defaultBackoff
		backoffStr := os.Getenv(envBackoff)
		if backoffStr == "" {
			backoff, err := time.ParseDuration(backoffStr)
			if err != nil {
				o.MaxBackoff = backoff
			}
		}
	})

	// Create SSM client with retryer
	client = ssm.NewFromConfig(cfg, func(o *ssm.Options) {
		o.Retryer = customRetryer
	})

	return nil
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

func fetchParams(paths []string) ([]types.Parameter, error) {
	// create tag filters
	tagFilters := make([]types.ParameterStringFilter, len(tags))
	for i, tag := range tags {
		tagFilter := fmt.Sprintf("tag:%s", tag)
		tagFilters[i] = types.ParameterStringFilter{
			Key: &tagFilter,
		}
	}

	// TEMP: until parameter-filters work for get-parameters-by-path
	// - https://docs.aws.amazon.com/cli/latest/reference/ssm/get-parameters-by-path.html ("This API action doesn't support filtering by tags.")
	// - https://github.com/aws/aws-cli/issues/2850)
	//
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

func describeParams(filters []types.ParameterStringFilter) ([]string, error) {
	paramNames := make([]string, 0)

	done := false
	var nextToken string
	for !done {
		input := &ssm.DescribeParametersInput{
			ParameterFilters: filters,
		}

		if nextToken != "" {
			input.NextToken = &nextToken
		}

		output, err := client.DescribeParameters(context.TODO(), input)
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

func getParamsByPath(paths []string) ([]types.Parameter, error) {
	params := make([]types.Parameter, 0)

	// retrieve params for all paths
	for _, path := range paths {

		done := false
		var nextToken string
		for !done {
			input := &ssm.GetParametersByPathInput{
				Path:           &path,
				Recursive:      true,
				WithDecryption: true,
			}

			if nextToken != "" {
				input.NextToken = &nextToken
			}

			output, err := client.GetParametersByPath(context.TODO(), input)
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

func calcUnion(paramNames []string, params []types.Parameter) []types.Parameter {
	// build map lookup
	lookup := make(map[string]bool)
	for _, paramName := range paramNames {
		lookup[paramName] = true
	}

	// calculate union
	union := make([]types.Parameter, 0)
	for _, param := range params {
		if lookup[*param.Name] {
			union = append(union, param)
		}
	}

	return union
}

func printParams(params []types.Parameter) {
	for _, param := range params {
		split := strings.Split(*param.Name, "/")
		name := split[len(split)-1]
		fmt.Printf("%s=%s\n", strings.ToUpper(name), *param.Value)
	}
}
