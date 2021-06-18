// Code generated by smithy-go-codegen DO NOT EDIT.

package ssm

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// A State Manager association defines the state that you want to maintain on your
// instances. For example, an association can specify that anti-virus software must
// be installed and running on your instances, or that certain ports must be
// closed. For static targets, the association specifies a schedule for when the
// configuration is reapplied. For dynamic targets, such as an AWS Resource Group
// or an AWS Autoscaling Group, State Manager applies the configuration when new
// instances are added to the group. The association also specifies actions to take
// when applying the configuration. For example, an association for anti-virus
// software might run once a day. If the software is not installed, then State
// Manager installs it. If the software is installed, but the service is not
// running, then the association might instruct State Manager to start the service.
func (c *Client) CreateAssociation(ctx context.Context, params *CreateAssociationInput, optFns ...func(*Options)) (*CreateAssociationOutput, error) {
	if params == nil {
		params = &CreateAssociationInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateAssociation", params, optFns, addOperationCreateAssociationMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateAssociationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateAssociationInput struct {

	// The name of the SSM document that contains the configuration information for the
	// instance. You can specify Command or Automation documents. You can specify
	// AWS-predefined documents, documents you created, or a document that is shared
	// with you from another account. For SSM documents that are shared with you from
	// other AWS accounts, you must specify the complete SSM document ARN, in the
	// following format: arn:partition:ssm:region:account-id:document/document-name
	// For example: arn:aws:ssm:us-east-2:12345678912:document/My-Shared-Document For
	// AWS-predefined documents and SSM documents you created in your account, you only
	// need to specify the document name. For example, AWS-ApplyPatchBaseline or
	// My-Document.
	//
	// This member is required.
	Name *string

	// By default, when you create a new association, the system runs it immediately
	// after it is created and then according to the schedule you specified. Specify
	// this option if you don't want an association to run immediately after you create
	// it. This parameter is not supported for rate expressions.
	ApplyOnlyAtCronInterval bool

	// Specify a descriptive name for the association.
	AssociationName *string

	// Specify the target for the association. This target is required for associations
	// that use an Automation document and target resources by using rate controls.
	AutomationTargetParameterName *string

	// The names or Amazon Resource Names (ARNs) of the Systems Manager Change Calendar
	// type documents you want to gate your associations under. The associations only
	// run when that Change Calendar is open. For more information, see AWS Systems
	// Manager Change Calendar
	// (https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-change-calendar).
	CalendarNames []string

	// The severity level to assign to the association.
	ComplianceSeverity types.AssociationComplianceSeverity

	// The document version you want to associate with the target(s). Can be a specific
	// version or the default version.
	DocumentVersion *string

	// The instance ID. InstanceId has been deprecated. To specify an instance ID for
	// an association, use the Targets parameter. Requests that include the parameter
	// InstanceID with SSM documents that use schema version 2.0 or later will fail. In
	// addition, if you use the parameter InstanceId, you cannot use the parameters
	// AssociationName, DocumentVersion, MaxErrors, MaxConcurrency, OutputLocation, or
	// ScheduleExpression. To use these parameters, you must use the Targets parameter.
	InstanceId *string

	// The maximum number of targets allowed to run the association at the same time.
	// You can specify a number, for example 10, or a percentage of the target set, for
	// example 10%. The default value is 100%, which means all targets run the
	// association at the same time. If a new instance starts and attempts to run an
	// association while Systems Manager is running MaxConcurrency associations, the
	// association is allowed to run. During the next association interval, the new
	// instance will process its association within the limit specified for
	// MaxConcurrency.
	MaxConcurrency *string

	// The number of errors that are allowed before the system stops sending requests
	// to run the association on additional targets. You can specify either an absolute
	// number of errors, for example 10, or a percentage of the target set, for example
	// 10%. If you specify 3, for example, the system stops sending requests when the
	// fourth error is received. If you specify 0, then the system stops sending
	// requests after the first error is returned. If you run an association on 50
	// instances and set MaxError to 10%, then the system stops sending the request
	// when the sixth error is received. Executions that are already running an
	// association when MaxErrors is reached are allowed to complete, but some of these
	// executions may fail as well. If you need to ensure that there won't be more than
	// max-errors failed executions, set MaxConcurrency to 1 so that executions proceed
	// one at a time.
	MaxErrors *string

	// An S3 bucket where you want to store the output details of the request.
	OutputLocation *types.InstanceAssociationOutputLocation

	// The parameters for the runtime configuration of the document.
	Parameters map[string][]string

	// A cron expression when the association will be applied to the target(s).
	ScheduleExpression *string

	// The mode for generating association compliance. You can specify AUTO or MANUAL.
	// In AUTO mode, the system uses the status of the association execution to
	// determine the compliance status. If the association execution runs successfully,
	// then the association is COMPLIANT. If the association execution doesn't run
	// successfully, the association is NON-COMPLIANT. In MANUAL mode, you must specify
	// the AssociationId as a parameter for the PutComplianceItems API action. In this
	// case, compliance data is not managed by State Manager. It is managed by your
	// direct call to the PutComplianceItems API action. By default, all associations
	// use AUTO mode.
	SyncCompliance types.AssociationSyncCompliance

	// A location is a combination of AWS Regions and AWS accounts where you want to
	// run the association. Use this action to create an association in multiple
	// Regions and multiple accounts.
	TargetLocations []types.TargetLocation

	// The targets for the association. You can target instances by using tags, AWS
	// Resource Groups, all instances in an AWS account, or individual instance IDs.
	// For more information about choosing targets for an association, see Using
	// targets and rate controls with State Manager associations
	// (https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-state-manager-targets-and-rate-controls.html)
	// in the AWS Systems Manager User Guide.
	Targets []types.Target
}

type CreateAssociationOutput struct {

	// Information about the association.
	AssociationDescription *types.AssociationDescription

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addOperationCreateAssociationMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpCreateAssociation{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpCreateAssociation{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpCreateAssociationValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateAssociation(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opCreateAssociation(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ssm",
		OperationName: "CreateAssociation",
	}
}
