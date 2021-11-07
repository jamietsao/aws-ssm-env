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

	fetchByPaths := len(paths) > 0
	fetchByTags := len(tags) > 0

	// if only fetching params by tags
	if fetchByTags && !fetchByPaths {
		// describe all parameters with given tags
		paramNames, err := describeParameters(tags)
		if err != nil {
			return nil, err
		}

		// fetch these parameters directly
		params, err := getParameters(paramNames)
		if err != nil {
			return params, err
		}
		return params, nil

	} else if fetchByPaths && !fetchByTags {
		// if only fetching params by paths

		// retrieve params for given paths
		params, err := getParametersByPath(paths)
		if err != nil {
			return params, err
		}

		return params, nil
	} else {
		// else if fetching params by both paths and tags

		// describe all parameters with given tags
		paramNames, err := describeParameters(tags)
		if err != nil {
			return nil, err
		}

		// retrieve params for given paths
		params, err := getParametersByPath(paths)
		if err != nil {
			return params, err
		}

		// calculate union of two sets
		union := calcUnion(paramNames, params)

		return union, nil
	}
}

// Gets information about SSM parameters with the given tags via describe-parameters (https://docs.aws.amazon.com/cli/latest/reference/ssm/describe-parameters.html)
func describeParameters(tags []string) ([]*string, error) {

	// create tag filters
	tagFilters := make([]*ssm.ParameterStringFilter, len(tags))
	for i, tag := range tags {
		tagFilter := fmt.Sprintf("tag:%s", tag)
		tagFilters[i] = &ssm.ParameterStringFilter{
			Key: &tagFilter,
		}
	}

	paramNames := make([]*string, 0)

	done := false
	var nextToken string
	for !done {
		input := &ssm.DescribeParametersInput{
			ParameterFilters: tagFilters,
		}

		if nextToken != "" {
			input.SetNextToken(nextToken)
		}

		output, err := client.DescribeParameters(input)
		if err != nil {
			return paramNames, err
		}

		for _, param := range output.Parameters {
			paramNames = append(paramNames, param.Name)
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

// Retrieves SSM parameters with the given names via get-parameters (https://docs.aws.amazon.com/cli/latest/reference/ssm/get-parameters.html)
func getParameters(paramNames []*string) ([]*ssm.Parameter, error) {

	// GetParameters only supports at max of 10 params
	chunks := chunkParamNames(paramNames, 10)

	parameters := make([]*ssm.Parameter, 0)
	for _, chunk := range chunks {
		output, err := client.GetParameters(&ssm.GetParametersInput{
			Names:          chunk,
			WithDecryption: &trueBool,
		})

		if err != nil {
			return nil, err
		}

		if len(output.InvalidParameters) > 1 {
			return nil, fmt.Errorf("Invalid parameters found %v", output.InvalidParameters)
		}

		parameters = append(parameters, output.Parameters...)
	}

	return parameters, nil
}

func chunkParamNames(paramNames []*string, chunkSize int) [][]*string {
	var chunks [][]*string
	for i := 0; i < len(paramNames); i += chunkSize {
		end := i + chunkSize
		if end > len(paramNames) {
			end = len(paramNames)
		}

		chunks = append(chunks, paramNames[i:end])
	}

	return chunks
}

// Retrieves SSM parameters in the given path hierarchies via get-parameters-by-path (https://docs.aws.amazon.com/cli/latest/reference/ssm/get-parameters-by-path.html)
func getParametersByPath(paths []string) ([]*ssm.Parameter, error) {
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

func calcUnion(paramNames []*string, params []*ssm.Parameter) []*ssm.Parameter {
	// build map lookup
	lookup := make(map[string]bool)
	for _, paramName := range paramNames {
		lookup[*paramName] = true
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
