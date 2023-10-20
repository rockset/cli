package format

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/rockset/rockset-go-client/openapi"
	"time"
)

var FieldFormatters = map[string]FieldFormatter{
	SizeName:        &SizeFormatter{},
	TimeSince:       &TimeSinceFormatter{},
	IntegrationType: &IntegrationTypeFormatter{},
}

type FieldFormatter interface {
	FormatField(any) (string, error)
	Name() string
}

type SizeFormatter struct{}

const SizeName = "size"

func (f SizeFormatter) Name() string { return SizeName }
func (f SizeFormatter) FormatField(a any) (string, error) {
	switch a.(type) {
	case int64:
		return humanize.Bytes(uint64(a.(int64))), nil
	case uint64:
		return humanize.Bytes(a.(uint64)), nil
	default:
		return "", fmt.Errorf("%v can't be turned into a size", a)
	}
}

const TimeSince = "time_since"

type TimeSinceFormatter struct{}

func (f TimeSinceFormatter) Name() string { return TimeSince }
func (f TimeSinceFormatter) FormatField(a any) (string, error) {
	switch a.(type) {
	case int64:
		i := a.(int64)
		if i == 0 {
			return "never", nil
		}
		t := time.UnixMilli(i)
		return humanize.Time(t), nil
	default:
		return "", fmt.Errorf("%v (%T) can't be turned into time", a, a)
	}
}

const IntegrationType = "integration_type"

type IntegrationTypeFormatter struct{}

func (f IntegrationTypeFormatter) Name() string { return IntegrationType }
func (f IntegrationTypeFormatter) FormatField(a any) (string, error) {
	switch integration := a.(type) {
	case openapi.Integration:
		if integration.AzureBlobStorage != nil {
			return "Azure Blob Storage", nil
		}
		if integration.AzureServiceBus != nil {
			return "Azure Service Bus", nil
		}
		if integration.AzureEventHubs != nil {
			return "Azure Event Hubs", nil
		}
		if integration.Dynamodb != nil {
			return "Amazon DynamoDB", nil
		}
		if integration.Gcs != nil {
			return "Google Cloud Storage", nil
		}
		if integration.Kafka != nil {
			if integration.Kafka.HasAwsRole() {
				return "Amazon MSK", nil
			} else if integration.Kafka.HasUseV3() {
				return "Confluent Cloud", nil
			} else {
				return "Kafka", nil
			}
		}
		if integration.Mongodb != nil {
			return "MongoDB", nil
		}
		if integration.Kinesis != nil {
			return "Amazon Kinesis", nil
		}
		if integration.S3 != nil {
			return "Amazon S3", nil
		}
		if integration.Snowflake != nil {
			return "Snowflake", nil
		}
		return "Unknown", nil
	default:
		return "", fmt.Errorf("can't parse integration type from non-Integration %v (%T)", a, a)
	}
}
