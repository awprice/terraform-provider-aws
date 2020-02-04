package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAwsEc2TransitGatewayVpcAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsEc2TransitGatewayVpcAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAwsEc2TransitGatewayVpcAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	input := &ec2.DescribeTransitGatewayVpcAttachmentsInput{}

	if v, ok := d.GetOk("filter"); ok {
		input.Filters = buildAwsDataSourceFilters(v.(*schema.Set))
	}

	transitGatewayVPCAttachmentsIDs := make([]string, 0)
	log.Printf("[DEBUG] Reading EC2 Transit Gateways VPC Attachments: %s", input)
	err := conn.DescribeTransitGatewayVpcAttachmentsPages(input, func(page *ec2.DescribeTransitGatewayVpcAttachmentsOutput, lastPage bool) bool {
		for _, transitGatewayVPCAttachment := range page.TransitGatewayVpcAttachments {
			transitGatewayVPCAttachmentsIDs = append(transitGatewayVPCAttachmentsIDs, aws.StringValue(transitGatewayVPCAttachment.TransitGatewayAttachmentId))
		}
		return !lastPage
	})
	if err != nil {
		return fmt.Errorf("error reading EC2 Transit Gateway VPC Attachments: %s", err)
	}

	d.SetId(resource.UniqueId())
	if err := d.Set("ids", transitGatewayVPCAttachmentsIDs); err != nil {
		return fmt.Errorf("error setting ids: %s", err)
	}

	return nil
}
