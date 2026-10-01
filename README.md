# Pulumi GCP Go: Pub/Sub

This project provisions GCP Pub/Sub resources using Pulumi and Go. It demonstrates how to:
- Use the Pulumi GCP provider in a Go program
- Create Pub/Sub topics and subscriptions
- Configure message retention, retry policies, and ordering

## Providers

- Google Cloud Platform via the Pulumi GCP SDK for Go (`github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp`)

## Resources

- **Project Service** (`gcp.projects.Service`)
  - Enables Pub/Sub API

- **Topic** (`gcp.pubsub.Topic`)
  - Pub/Sub topic with configurable retention and labels

- **Subscription** (`gcp.pubsub.Subscription`)
  - Pull or push subscription with ack deadlines, retry policies, and expiration

## Outputs

- **topicId**: ID of the Pub/Sub topic
- **topicName**: Name of the Pub/Sub topic
- **subscriptionId**: ID of the Pub/Sub subscription
- **subscriptionName**: Name of the Pub/Sub subscription

## When to Use This Project

- You need a message queue for async communication between services
- You want to decouple event producers from consumers
- You need reliable message delivery with configurable retry policies

## Prerequisites

- Go 1.21+ installed
- A Google Cloud account with billing enabled
- GCP credentials configured for Pulumi (via `gcloud auth application-default login`)

## Usage

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Configure your stack:
   ```bash
   cp Pulumi.dev.yaml.example Pulumi.dev.yaml
   pulumi config set gcp:project YOUR_PROJECT_ID
   pulumi config set pubsub:topicName my-topic
   pulumi config set pubsub:subscriptionName my-subscription
   ```

3. Preview and deploy:
   ```bash
   pulumi preview
   pulumi up
   ```

## Project Layout

```
├── Pulumi.yaml                  Pulumi project definition
├── Pulumi.dev.yaml.example      Template for local dev configuration
├── go.mod                       Go module declaration and dependencies
├── main.go                      Pulumi program defining Pub/Sub resources
├── .gitignore                   Git ignore rules
└── LICENSE                      MIT License
```

## Configuration

| Name | Description | Default |
|------|-------------|---------|
| `gcp:project` | The Google Cloud project to deploy into | _required_ |
| `gcp:region` | The GCP region for resources | `us-central1` |
| `pubsub:topicName` | Name of the Pub/Sub topic | `my-topic` |
| `pubsub:messageRetentionDuration` | Message retention duration | `86600s` |
| `pubsub:labels` | Labels for the topic (JSON) | `{"environment": "dev"}` |
| `pubsub:subscriptionName` | Name of the subscription | `my-subscription` |
| `pubsub:ackDeadlineSeconds` | Acknowledgement deadline | `20` |
| `pubsub:pushEndpoint` | Push endpoint URL (empty for pull) | `""` |
| `pubsub:pushAttributes` | Push attributes (JSON) | `{}` |
| `pubsub:expirationTtl` | Subscription expiration TTL | `86400s` |
| `pubsub:minRetryBackoff` | Minimum retry backoff | `10s` |
| `pubsub:maxRetryBackoff` | Maximum retry backoff | `600s` |
| `pubsub:enableMessageOrdering` | Enable message ordering | `false` |

## Next Steps

- Add dead letter topics for failed message handling
- Configure schema validation for messages
- Set up message filtering for subscriptions
- Integrate with Cloud Functions or Cloud Run for event processing

## Getting Help

- Pulumi Documentation: https://www.pulumi.com/docs/
- GCP Provider Reference: https://www.pulumi.com/registry/packages/gcp/
- Community Slack: https://slack.pulumi.com/
- GitHub Issues: https://github.com/pulumi/pulumi/issues