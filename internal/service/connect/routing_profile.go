package connect

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func ResourceRoutingProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceRoutingProfileRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_outbound_queue_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 250),
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"media_concurrencies": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(connect.Channel_Values(), false), // Valid values: VOICE | CHAT | TASK
						},
						"concurrency": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 127),
			},
			"queue_configs": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(connect.Channel_Values(), false), // Valid values: VOICE | CHAT | TASK
						},
						"delay": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 9999),
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 99),
						},
						"queue_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"queue_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"routing_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},
	}
}

func resourceRoutingProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	instanceID, routingProfileID, err := RoutingProfileParseID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.DescribeRoutingProfileWithContext(ctx, &connect.DescribeRoutingProfileInput{
		InstanceId:       aws.String(instanceID),
		RoutingProfileId: aws.String(routingProfileID),
	})

	if !d.IsNewResource() && tfawserr.ErrMessageContains(err, connect.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] Connect Routing Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Connect Routing Profile (%s): %w", d.Id(), err))
	}

	if resp == nil || resp.RoutingProfile == nil {
		return diag.FromErr(fmt.Errorf("error getting Connect Routing Profile (%s): empty response", d.Id()))
	}

	routingProfile := resp.RoutingProfile

	if err := d.Set("media_concurrencies", flattenRoutingProfileMediaConcurrencies(routingProfile.MediaConcurrencies)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("arn", routingProfile.RoutingProfileArn)
	d.Set("default_outbound_queue_id", routingProfile.DefaultOutboundQueueId)
	d.Set("description", routingProfile.Description)
	d.Set("instance_id", instanceID)
	d.Set("name", routingProfile.Name)

	d.Set("routing_profile_id", routingProfile.RoutingProfileId)

	// getting the routing profile queues uses a separate API: ListRoutingProfileQueues
	queueConfigs, err := getConnectRoutingProfileQueueConfigs(ctx, conn, instanceID, routingProfileID)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error finding Connect Routing Profile Queue Configs Summary by Routing Profile ID (%s): %w", routingProfileID, err))
	}

	d.Set("queue_configs", queueConfigs)

	tags := KeyValueTags(routingProfile.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tags_all: %w", err))
	}

	return nil
}

func flattenRoutingProfileMediaConcurrencies(mediaConcurrencies []*connect.MediaConcurrency) []interface{} {
	mediaConcurrenciesList := []interface{}{}

	for _, mediaConcurrency := range mediaConcurrencies {
		values := map[string]interface{}{
			"channel":     aws.StringValue(mediaConcurrency.Channel),
			"concurrency": aws.Int64Value(mediaConcurrency.Concurrency),
		}

		mediaConcurrenciesList = append(mediaConcurrenciesList, values)
	}
	return mediaConcurrenciesList
}

func getConnectRoutingProfileQueueConfigs(ctx context.Context, conn *connect.Connect, instanceID, routingProfileID string) ([]interface{}, error) {
	queueConfigsList := []interface{}{}

	input := &connect.ListRoutingProfileQueuesInput{
		InstanceId:       aws.String(instanceID),
		MaxResults:       aws.Int64(ListRoutingProfileQueuesMaxResults),
		RoutingProfileId: aws.String(routingProfileID),
	}

	err := conn.ListRoutingProfileQueuesPagesWithContext(ctx, input, func(page *connect.ListRoutingProfileQueuesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, qc := range page.RoutingProfileQueueConfigSummaryList {
			if qc == nil {
				continue
			}

			values := map[string]interface{}{
				"channel":    aws.StringValue(qc.Channel),
				"delay":      aws.Int64Value(qc.Delay),
				"priority":   aws.Int64Value(qc.Priority),
				"queue_arn":  aws.StringValue(qc.QueueArn),
				"queue_id":   aws.StringValue(qc.QueueId),
				"queue_name": aws.StringValue(qc.QueueName),
			}

			queueConfigsList = append(queueConfigsList, values)
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return queueConfigsList, nil
}

func RoutingProfileParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected instanceID:routingProfileID", id)
	}

	return parts[0], parts[1], nil
}
