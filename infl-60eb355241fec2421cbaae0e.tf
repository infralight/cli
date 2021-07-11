resource "aws_s3_bucket" "aws-cloudtrail-logs-094724549126-0b9a848c" {
	bucket = "aws-cloudtrail-logs-094724549126-0b9a848c"
	force_destroy = false
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "IPAllow",
			"Effect": "Deny",
			"Principal": "*",
			"Action": [
				"s3:PutObject",
				"s3:GetObject"
			],
			"Resource": [
				"arn:aws:s3:::aws-cloudtrail-logs-094724549126-0b9a848c",
				"arn:aws:s3:::aws-cloudtrail-logs-094724549126-0b9a848c/*"
			],
			"Condition": {
				"Bool": {
					"aws:SecureTransport": "false"
				}
			}
		},
		{
			"Sid": "AWSCloudTrailAclCheck20150319",
			"Effect": "Allow",
			"Principal": {
				"Service": "cloudtrail.amazonaws.com"
			},
			"Action": "s3:GetBucketAcl",
			"Resource": "arn:aws:s3:::aws-cloudtrail-logs-094724549126-0b9a848c"
		},
		{
			"Sid": "AWSCloudTrailWrite20150319",
			"Effect": "Allow",
			"Principal": {
				"Service": "cloudtrail.amazonaws.com"
			},
			"Action": "s3:PutObject",
			"Resource": "arn:aws:s3:::aws-cloudtrail-logs-094724549126-0b9a848c/null/AWSLogs/094724549126/*",
			"Condition": {
				"StringEquals": {
					"s3:x-amz-acl": "bucket-owner-full-control"
				}
			}
		}
	]
}
POLICY
	server_side_encryption_configuration {
		rule {
			apply_server_side_encryption_by_default {
				sse_algorithm = "aws:kms"
				kms_master_key_id = "arn:aws:kms:us-east-1:094724549126:key/8fc17849-d564-49fd-9347-f62641e2c32e"
			}
		}
	}
	logging {
		target_bucket = "prod-infralight-log-bucket"
		target_prefix = "log/cloudtrail_logs/"
	}
	versioning {
		enabled = true
		mfa_delete = true
	}
}