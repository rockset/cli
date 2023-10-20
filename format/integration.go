package format

import "github.com/rockset/rockset-go-client/openapi"

var IntegrationDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Name", "name"),
		{
			ColumnName: "Collections",
			Path: []PathElem{{
				FieldName:       "collections",
				HasArrayMapping: true,
			}, {
				FieldName: "name",
			}},
		},
		// TODO need a selector which shows the type
	},
	Wide: []FieldSelection{
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Description", "description"),
		NewFieldSelection("Created By", "created_by"),
		NewFieldSelection("Created At", "created_at"),
	},
}

// just to list available fields
var _ = openapi.Integration{
	AzureBlobStorage: &openapi.AzureBlobStorageIntegration{
		ConnectionString: "",
	},
	AzureEventHubs: &openapi.AzureEventHubsIntegration{
		ConnectionString: nil,
	},
	AzureServiceBus: &openapi.AzureServiceBusIntegration{
		ConnectionString: "",
	},
	Collections:         nil,
	CreatedAt:           nil,
	CreatedBy:           "",
	CreatedByApikeyName: nil,
	Description:         nil,
	Dynamodb: &openapi.DynamodbIntegration{
		AwsAccessKey: &openapi.AwsAccessKey{
			AwsAccessKeyId:     "",
			AwsSecretAccessKey: "",
		},
		AwsRole: &openapi.AwsRole{
			AwsExternalId: nil,
			AwsRoleArn:    "",
		},
		S3ExportBucketName: nil,
	},
	Gcs: &openapi.GcsIntegration{
		GcpServiceAccount: &openapi.GcpServiceAccount{
			ServiceAccountKeyFileJson: "",
		},
	},
	Kafka: &openapi.KafkaIntegration{
		AwsRole: &openapi.AwsRole{
			AwsExternalId: nil,
			AwsRoleArn:    "",
		},
		BootstrapServers: nil,
		ConnectionString: nil,
		KafkaDataFormat:  nil,
		KafkaTopicNames:  nil,
		SchemaRegistryConfig: &openapi.SchemaRegistryConfig{
			Key:    nil,
			Secret: nil,
			Url:    nil,
		},
		SecurityConfig: &openapi.KafkaV3SecurityConfig{
			ApiKey: nil,
			Secret: nil,
		},
		SourceStatusByTopic: nil,
		UseV3:               nil,
	},
	Kinesis: &openapi.KinesisIntegration{
		AwsAccessKey: &openapi.AwsAccessKey{
			AwsAccessKeyId:     "",
			AwsSecretAccessKey: "",
		},
		AwsRole: &openapi.AwsRole{
			AwsExternalId: nil,
			AwsRoleArn:    "",
		},
	},
	Mongodb: &openapi.MongoDbIntegration{
		ConnectionUri: "",
		Tls: &openapi.TLSConfig{
			CaCert:            nil,
			ClientCert:        "",
			ClientCertExpiry:  nil,
			ClientCertSubject: nil,
			ClientKey:         "",
		},
	},
	Name:       "",
	OwnerEmail: nil,
	S3: &openapi.S3Integration{
		AwsAccessKey: &openapi.AwsAccessKey{
			AwsAccessKeyId:     "",
			AwsSecretAccessKey: "",
		},
		AwsRole: &openapi.AwsRole{
			AwsExternalId: nil,
			AwsRoleArn:    "",
		},
	},
	Snowflake: &openapi.SnowflakeIntegration{
		AwsAccessKey: &openapi.AwsAccessKey{
			AwsAccessKeyId:     "",
			AwsSecretAccessKey: "",
		},
		AwsRole: &openapi.AwsRole{
			AwsExternalId: nil,
			AwsRoleArn:    "",
		},
		DefaultWarehouse: "",
		Password:         "",
		S3ExportPath:     "",
		SnowflakeUrl:     "",
		UserRole:         nil,
		Username:         "",
	},
}
