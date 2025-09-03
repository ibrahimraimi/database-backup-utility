# Cloud Storage Integration Guide

This guide covers integrating the Database Backup Utility with various cloud storage providers.

## Supported Cloud Providers

- **AWS S3** - Amazon Simple Storage Service
- **Google Cloud Storage** - Google Cloud Platform storage
- **Azure Blob Storage** - Microsoft Azure storage

## AWS S3 Integration

### Prerequisites

- AWS account with S3 access
- IAM user or role with appropriate permissions
- S3 bucket created for backups

### IAM Permissions

Create an IAM policy with the following permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::your-backup-bucket",
        "arn:aws:s3:::your-backup-bucket/*"
      ]
    }
  ]
}
```

### Configuration

```yaml
# ~/dbu.yaml
storage:
  type: "cloud"
  path: "./temp"

cloud:
  provider: "aws"
  bucket: "your-backup-bucket"
  region: "us-east-1"
  access_key: "" # Set via AWS_ACCESS_KEY_ID
  secret_key: "" # Set via AWS_SECRET_ACCESS_KEY
```

### Environment Variables

```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-east-1
export AWS_S3_BUCKET=your-backup-bucket
```

### Usage Examples

#### Backup to S3

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --storage cloud \
  --cloud-provider aws \
  --bucket your-backup-bucket \
  --region us-east-1 \
  --compress
```

#### Restore from S3

```bash
./dbu restore \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --file s3://your-backup-bucket/mysql_mydb_full_2024-01-15_10-30-00.sql.gz
```

### S3-Specific Features

#### Server-Side Encryption

```bash
# Backup with server-side encryption
./dbu backup \
  --db-type mysql \
  --storage cloud \
  --cloud-provider aws \
  --bucket your-backup-bucket \
  --compress
```

#### S3 Transfer Acceleration

```bash
# Enable transfer acceleration in AWS CLI
aws configure set default.s3.use_accelerate_endpoint true
```

#### S3 Lifecycle Policies

```json
{
  "Rules": [
    {
      "ID": "BackupLifecycle",
      "Status": "Enabled",
      "Transitions": [
        {
          "Days": 30,
          "StorageClass": "STANDARD_IA"
        },
        {
          "Days": 90,
          "StorageClass": "GLACIER"
        }
      ],
      "Expiration": {
        "Days": 2555
      }
    }
  ]
}
```

## Google Cloud Storage Integration

### Prerequisites

- Google Cloud Platform account
- Service account with Storage permissions
- GCS bucket created for backups

### Service Account Setup

1. Create a service account in Google Cloud Console
2. Download the service account key file
3. Grant the following roles:
   - Storage Object Admin
   - Storage Object Viewer

### Configuration

```yaml
# ~/dbu.yaml
storage:
  type: "cloud"
  path: "./temp"

cloud:
  provider: "gcp"
  bucket: "your-backup-bucket"
  region: "us-central1"
  access_key: "" # Set via GOOGLE_APPLICATION_CREDENTIALS
  secret_key: ""
```

### Environment Variables

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
export GCP_PROJECT_ID=your-project-id
export GCP_BUCKET=your-backup-bucket
```

### Usage Examples

#### Backup to GCS

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --storage cloud \
  --cloud-provider gcp \
  --bucket your-backup-bucket \
  --region us-central1 \
  --compress
```

#### Restore from GCS

```bash
./dbu restore \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --file gs://your-backup-bucket/mysql_mydb_full_2024-01-15_10-30-00.sql.gz
```

### GCS-Specific Features

#### Storage Classes

```bash
# Use different storage classes for cost optimization
# Nearline: 30-day minimum storage
# Coldline: 90-day minimum storage
# Archive: 365-day minimum storage
```

#### Object Lifecycle Management

```json
{
  "lifecycle": {
    "rule": [
      {
        "action": {
          "type": "SetStorageClass",
          "storageClass": "NEARLINE"
        },
        "condition": {
          "age": 30
        }
      },
      {
        "action": {
          "type": "SetStorageClass",
          "storageClass": "COLDLINE"
        },
        "condition": {
          "age": 90
        }
      },
      {
        "action": {
          "type": "Delete"
        },
        "condition": {
          "age": 2555
        }
      }
    ]
  }
}
```

## Azure Blob Storage Integration

### Prerequisites

- Azure account with Blob Storage access
- Storage account created
- Container created for backups

### Storage Account Setup

1. Create a storage account in Azure Portal
2. Create a container for backups
3. Generate access keys or use managed identity

### Configuration

```yaml
# ~/dbu.yaml
storage:
  type: "cloud"
  path: "./temp"

cloud:
  provider: "azure"
  bucket: "your-backup-container"
  region: "eastus"
  access_key: "" # Set via AZURE_STORAGE_ACCOUNT
  secret_key: "" # Set via AZURE_STORAGE_KEY
```

### Environment Variables

```bash
export AZURE_STORAGE_ACCOUNT=your_storage_account
export AZURE_STORAGE_KEY=your_storage_key
export AZURE_CONTAINER=your-backup-container
```

### Usage Examples

#### Backup to Azure Blob

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --storage cloud \
  --cloud-provider azure \
  --bucket your-backup-container \
  --region eastus \
  --compress
```

#### Restore from Azure Blob

```bash
./dbu restore \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --file azure://your-backup-container/mysql_mydb_full_2024-01-15_10-30-00.sql.gz
```

### Azure-Specific Features

#### Access Tiers

```bash
# Use different access tiers for cost optimization
# Hot: Frequently accessed data
# Cool: Infrequently accessed data
# Archive: Rarely accessed data
```

#### Blob Lifecycle Management

```json
{
  "rules": [
    {
      "name": "BackupLifecycle",
      "enabled": true,
      "type": "Lifecycle",
      "definition": {
        "filters": {
          "blobTypes": ["blockBlob"]
        },
        "actions": {
          "baseBlob": {
            "tierToCool": {
              "daysAfterModificationGreaterThan": 30
            },
            "tierToArchive": {
              "daysAfterModificationGreaterThan": 90
            },
            "delete": {
              "daysAfterModificationGreaterThan": 2555
            }
          }
        }
      }
    }
  ]
}
```

## Multi-Cloud Backup Strategy

### Backup to Multiple Providers

```bash
#!/bin/bash
# multi-cloud-backup.sh

DB_HOST="localhost"
DB_USER="root"
DB_PASS="mypassword"
DB_NAME="mydb"

# Backup to AWS S3
./dbu backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --storage cloud \
  --cloud-provider aws \
  --bucket aws-backup-bucket \
  --region us-east-1 \
  --compress

# Backup to Google Cloud Storage
./dbu backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --storage cloud \
  --cloud-provider gcp \
  --bucket gcp-backup-bucket \
  --region us-central1 \
  --compress

# Backup to Azure Blob Storage
./dbu backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --storage cloud \
  --cloud-provider azure \
  --bucket azure-backup-container \
  --region eastus \
  --compress
```

## Cost Optimization

### Storage Class Optimization

| Provider | Storage Class | Use Case          | Cost   |
| -------- | ------------- | ----------------- | ------ |
| AWS S3   | Standard      | Frequent access   | High   |
| AWS S3   | Standard-IA   | Infrequent access | Medium |
| AWS S3   | Glacier       | Archive           | Low    |
| GCS      | Standard      | Frequent access   | High   |
| GCS      | Nearline      | Infrequent access | Medium |
| GCS      | Coldline      | Archive           | Low    |
| Azure    | Hot           | Frequent access   | High   |
| Azure    | Cool          | Infrequent access | Medium |
| Azure    | Archive       | Archive           | Low    |

### Lifecycle Policies

Implement lifecycle policies to automatically move backups to cheaper storage classes:

```bash
# AWS S3 lifecycle policy
aws s3api put-bucket-lifecycle-configuration \
  --bucket your-backup-bucket \
  --lifecycle-configuration file://lifecycle.json

# Google Cloud Storage lifecycle policy
gsutil lifecycle set lifecycle.json gs://your-backup-bucket

# Azure Blob Storage lifecycle policy
az storage blob service-properties update \
  --account-name your_storage_account \
  --lifecycle-policy lifecycle.json
```

## Security Best Practices

### Access Control

1. **Use IAM roles instead of access keys when possible**
2. **Implement least privilege access**
3. **Enable MFA for administrative access**
4. **Regularly rotate access keys**
5. **Use bucket policies to restrict access**

### Encryption

1. **Enable server-side encryption**
2. **Use customer-managed keys when possible**
3. **Encrypt data in transit**
4. **Consider client-side encryption for sensitive data**

### Monitoring

1. **Enable CloudTrail (AWS) or equivalent**
2. **Set up billing alerts**
3. **Monitor access patterns**
4. **Implement anomaly detection**

## Troubleshooting

### Common Issues

#### Authentication Failures

```bash
# AWS S3
aws sts get-caller-identity

# Google Cloud Storage
gcloud auth list

# Azure Blob Storage
az account show
```

#### Network Issues

```bash
# Test connectivity
ping s3.amazonaws.com
ping storage.googleapis.com
ping blob.core.windows.net

# Check DNS resolution
nslookup s3.amazonaws.com
nslookup storage.googleapis.com
nslookup blob.core.windows.net
```

#### Permission Issues

```bash
# Test S3 permissions
aws s3 ls s3://your-backup-bucket

# Test GCS permissions
gsutil ls gs://your-backup-bucket

# Test Azure permissions
az storage blob list --container-name your-backup-container
```

### Performance Optimization

#### Upload Optimization

```bash
# Use multipart uploads for large files
# Enable transfer acceleration (AWS)
# Use parallel uploads
# Compress data before upload
```

#### Download Optimization

```bash
# Use CDN when possible
# Enable transfer acceleration (AWS)
# Use parallel downloads
# Cache frequently accessed data
```

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [MySQL Guide](mysql.md)
- [PostgreSQL Guide](postgresql.md)
- [MongoDB Guide](mongodb.md)
- [SQLite Guide](sqlite.md)
- [Troubleshooting](troubleshooting.md)
