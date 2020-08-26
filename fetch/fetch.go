package fetch

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	client *ssm.SSM
	paths  []string
	tags   []string

	trueBool = true
)

func init() {
	awsConfig := aws.NewConfig()
	ssmRegion := os.Getenv("SSM_REGION")
	if ssmRegion != "" {
		awsConfig = awsConfig.WithRegion(ssmRegion)
	}
	session := session.Must(session.NewSession())
	client = ssm.New(session, awsConfig)
}

func FetchParams(paths, tags []string) ([]*ssm.Parameter, error) {
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

func getParamNameValues(params []*ssm.Parameter) map[string]string {
	paramVals := make(map[string]string, len(params))
	for _, param := range params {
		split := strings.Split(*param.Name, "/")
		name := split[len(split)-1]
		paramVals[strings.ToUpper(name)] = *param.Value
	}
	return paramVals
}

func MustSetOS(paths, tags []string) {
	params, err := FetchParams(paths, tags)
	if err != nil {
		panic(err)
	}
	nameValues := getParamNameValues(params)
	for name, value := range nameValues {
		if err := os.Setenv(name, value); err != nil {
			panic(err)
		}
	}
}
