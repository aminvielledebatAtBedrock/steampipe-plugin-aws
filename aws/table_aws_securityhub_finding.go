package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHubFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_finding",
		Description: "AWS Security Hub Finding",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSecurityHubFinding,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubFindings,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "company_name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "compliance_status", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "confidence", Require: plugin.Optional, Operators: []string{"=", ">=", "<="}},
				{Name: "criticality", Require: plugin.Optional, Operators: []string{"=", ">=", "<="}},
				{Name: "generator_id", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "product_arn", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "product_name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "record_state", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "title", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "verification_state", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "workflow_state", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The security findings provider-specific identifier for a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "company_name",
				Description: "The name of the company for the product that generated the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "confidence",
				Description: "A finding's confidence. Confidence is defined as the likelihood that a finding accurately identifies the behavior or issue that it was intended to identify.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "created_at",
				Description: "Indicates when the security-findings provider created the potential security issue that a finding captured.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "compliance_status",
				Description: "The result of a compliance standards check.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Compliance.Status"),
			},
			{
				Name:        "updated_at",
				Description: "Indicates when the security-findings provider last updated the finding record.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "criticality",
				Description: "The level of importance assigned to the resources associated with the finding.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "A finding's description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "first_observed_at",
				Description: "Indicates when the security-findings provider first observed the potential security issue that a finding captured.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "generator_id",
				Description: "The identifier for the solution-specific component (a discrete unit of logic) that generated a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_observed_at",
				Description: "Indicates when the security-findings provider most recently observed the potential security issue that a finding captured.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "product_arn",
				Description: "The ARN generated by Security Hub that uniquely identifies a product that generates findings.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_name",
				Description: "The name of the product that generated the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "record_state",
				Description: "The record state of a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema_version",
				Description: "The schema version that a finding is formatted for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_url",
				Description: "A URL that links to a page about the current finding in the security-findings provider's solution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "verification_state",
				Description: "Indicates the veracity of a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workflow_state",
				Description: "The workflow state of a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "standards_control_arn",
				Description: "The ARN of the security standard control.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractStandardControlArn),
			},
			{
				Name:        "action",
				Description: "Provides details about an action that affects or that was taken on a resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "compliance",
				Description: "This data type is exclusive to findings that are generated as the result of a check run against a specific rule in a supported security standard, such as CIS Amazon Web Services Foundations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "finding_provider_fields",
				Description: "In a BatchImportFindings request, finding providers use FindingProviderFields to provide and update their own values for confidence, criticality, related findings, severity, and types.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "malware",
				Description: "A list of malware related to a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network",
				Description: "The details of network-related information about a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_path",
				Description: "Provides information about a network path that is relevant to a finding. Each entry under NetworkPath represents a component of that path.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "note",
				Description: "A user-defined note added to a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "patch_summary",
				Description: "Provides an overview of the patch compliance status for an instance against a selected compliance standard.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "process",
				Description: "The details of process-related information about a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product_fields",
				Description: "A data type where security-findings providers can include additional solution-specific details that aren't part of the defined AwsSecurityFinding format.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "related_findings",
				Description: "A list of related findings.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "remediation",
				Description: "A data type that describes the remediation options for a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resources",
				Description: "A set of resource data types that describe the resources that the finding refers to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "severity",
				Description: "A finding's severity.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "threat_intel_indicators",
				Description: "Threat intelligence details related to a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_defined_fields",
				Description: "A list of name/value string pairs associated with the finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vulnerabilities",
				Description: "Provides a list of vulnerabilities associated with the findings.",
				Type:        proto.ColumnType_JSON,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "A finding's title.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSecurityHubFindings")

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}
	input := &securityhub.GetFindingsInput{
		MaxResults: aws.Int64(100),
	}

	findingsFilter := buildListFindingsParam(d.Quals)
	if findingsFilter != nil {
		input.Filters = findingsFilter
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.GetFindingsPages(
		input,
		func(page *securityhub.GetFindingsOutput, isLast bool) bool {
			for _, finding := range page.Findings {
				d.StreamListItem(ctx, finding)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listSecurityHubFindings", "Error", err)
		// Handle error for unsupported or inactive regions
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHubFinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityHubFinding")

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Empty check
	if id == "" {
		return nil, nil
	}

	// get service
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &securityhub.GetFindingsInput{
		Filters: &securityhub.AwsSecurityFindingFilters{
			Id: []*securityhub.StringFilter{
				{
					Comparison: aws.String("EQUALS"),
					Value:      aws.String(id),
				},
			},
		},
	}

	// Get call
	op, err := svc.GetFindings(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getSecurityHubFinding", "ERROR", err)
		// Handle error for unsupported or inactive regions
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}

		return nil, err
	}
	if len(op.Findings) > 0 {
		return op.Findings[0], nil
	}
	return nil, nil
}

// Build param for findings list call
func buildListFindingsParam(quals plugin.KeyColumnQualMap) *securityhub.AwsSecurityFindingFilters {
	securityFindingsFilter := &securityhub.AwsSecurityFindingFilters{}
	strFilter := &securityhub.StringFilter{}

	strColumns := []string{"company_name", "compliance_status", "generator_id", "product_arn", "product_name", "record_state", "title", "verification_state", "workflow_state"}

	for _, s := range strColumns {
		if quals[s] == nil {
			continue
		}
		for _, q := range quals[s].Quals {
			value := q.Value.GetStringValue()
			if value == "" {
				continue
			}

			switch q.Operator {
			case "<>":
				strFilter.Comparison = aws.String("NOT_EQUALS")
			case "=":
				strFilter.Comparison = aws.String("EQUALS")
			}

			switch s {
			case "company_name":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.CompanyName = append(securityFindingsFilter.CompanyName, strFilter)
			case "generator_id":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.GeneratorId = append(securityFindingsFilter.GeneratorId, strFilter)
			case "compliance_status":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.ComplianceStatus = append(securityFindingsFilter.ComplianceStatus, strFilter)
			case "product_arn":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.ProductArn = append(securityFindingsFilter.ProductArn, strFilter)
			case "product_name":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.ProductName = append(securityFindingsFilter.ProductName, strFilter)
			case "record_state":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.RecordState = append(securityFindingsFilter.RecordState, strFilter)
			case "title":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.Title = append(securityFindingsFilter.Title, strFilter)
			case "verification_state":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.VerificationState = append(securityFindingsFilter.VerificationState, strFilter)
			case "workflow_state":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.WorkflowState = append(securityFindingsFilter.WorkflowState, strFilter)
			}

		}
	}

	return securityFindingsFilter
}

//// TRANSFORM FUNCTIONS

func extractStandardControlArn(_ context.Context, d *transform.TransformData) (interface{}, error) {
	findingArn := d.HydrateItem.(*securityhub.AwsSecurityFinding).Id

	if strings.Contains(*findingArn, "arn:aws:securityhub") {
		standardControlArn := strings.Replace(strings.Split(*findingArn, "/finding")[0], "subscription", "control", 1)
		return standardControlArn, nil
	}
	return nil, nil
}
