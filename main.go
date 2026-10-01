package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		projectID := config.Require(ctx, "gcp:project")
		region := config.Get(ctx, "gcp:region")
		if region == "" {
			region = "us-central1"
		}

		topicName := config.Get(ctx, "pubsub:topicName")
		if topicName == "" {
			topicName = "my-topic"
		}

		messageRetentionDuration := config.Get(ctx, "pubsub:messageRetentionDuration")
		if messageRetentionDuration == "" {
			messageRetentionDuration = "86600s"
		}

		labelsJSON := config.Get(ctx, "pubsub:labels")
		labels := pulumi.StringMap{}
		if labelsJSON != "" {
			// For simplicity, we'll use a fixed map
			labels = pulumi.StringMap{
				"environment": pulumi.String("dev"),
			}
		}

		subscriptionName := config.Get(ctx, "pubsub:subscriptionName")
		if subscriptionName == "" {
			subscriptionName = "my-subscription"
		}

		ackDeadlineSeconds := config.GetInt(ctx, "pubsub:ackDeadlineSeconds")
		if ackDeadlineSeconds == 0 {
			ackDeadlineSeconds = 20
		}

		pushEndpoint := config.Get(ctx, "pubsub:pushEndpoint")
		pushAttributesJSON := config.Get(ctx, "pubsub:pushAttributes")
		pushAttributes := pulumi.StringMap{}
		if pushAttributesJSON != "" {
			var attrs map[string]string
			if err := json.Unmarshal([]byte(pushAttributesJSON), &attrs); err != nil {
				return fmt.Errorf("invalid pubsub:pushAttributes JSON: %w", err)
			}
			pushAttributes = pulumi.ToStringMap(attrs)
		}

		expirationTtl := config.Get(ctx, "pubsub:expirationTtl")
		if expirationTtl == "" {
			expirationTtl = "86400s"
		}

		minRetryBackoff := config.Get(ctx, "pubsub:minRetryBackoff")
		if minRetryBackoff == "" {
			minRetryBackoff = "10s"
		}

		maxRetryBackoff := config.Get(ctx, "pubsub:maxRetryBackoff")
		if maxRetryBackoff == "" {
			maxRetryBackoff = "600s"
		}

		enableMessageOrdering := config.GetBool(ctx, "pubsub:enableMessageOrdering")

		// Enable Pub/Sub API
		pubsubAPI, err := projects.NewService(ctx, "pubsub-api", &projects.ServiceArgs{
			Service:          pulumi.String("pubsub.googleapis.com"),
			DisableOnDestroy: pulumi.Bool(false),
			Project:          pulumi.String(projectID),
		})
		if err != nil {
			return err
		}

		// Create Pub/Sub topic
		topic, err := pubsub.NewTopic(ctx, "topic", &pubsub.TopicArgs{
			Name:                     pulumi.String(topicName),
			MessageRetentionDuration: pulumi.String(messageRetentionDuration),
			Labels:                   labels,
		}, pulumi.DependsOn([]pulumi.Resource{pubsubAPI}))
		if err != nil {
			return err
		}

		// Create Pub/Sub subscription
		var pushConfig *pubsub.SubscriptionPushConfigArgs
		if pushEndpoint != "" {
			pushConfig = &pubsub.SubscriptionPushConfigArgs{
				PushEndpoint: pulumi.String(pushEndpoint),
				Attributes:   pushAttributes,
			}
		}

		subscription, err := pubsub.NewSubscription(ctx, "subscription", &pubsub.SubscriptionArgs{
			Name:               pulumi.String(subscriptionName),
			Topic:              topic.ID(),
			AckDeadlineSeconds: pulumi.Int(ackDeadlineSeconds),
			PushConfig:         pushConfig,
			ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
				Ttl: pulumi.String(expirationTtl),
			},
			RetryPolicy: &pubsub.SubscriptionRetryPolicyArgs{
				MinimumBackoff: pulumi.String(minRetryBackoff),
				MaximumBackoff: pulumi.String(maxRetryBackoff),
			},
			EnableMessageOrdering: pulumi.Bool(enableMessageOrdering),
		}, pulumi.DependsOn([]pulumi.Resource{topic}))
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("topicId", topic.ID())
		ctx.Export("topicName", topic.Name)
		ctx.Export("subscriptionId", subscription.ID())
		ctx.Export("subscriptionName", subscription.Name)

		return nil
	})
}
