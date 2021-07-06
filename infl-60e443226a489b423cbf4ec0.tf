resource "aws_s3_bucket" "config-bucket-094724549126" {
	bucket = "config-bucket-094724549126"
	force_destroy = false
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "AWSConfigBucketPermissionsCheck",
			"Effect": "Allow",
			"Principal": {
				"Service": "config.amazonaws.com"
			},
			"Action": "s3:GetBucketAcl",
			"Resource": "arn:aws:s3:::config-bucket-094724549126"
		},
		{
			"Sid": "AWSConfigBucketExistenceCheck",
			"Effect": "Allow",
			"Principal": {
				"Service": "config.amazonaws.com"
			},
			"Action": "s3:ListBucket",
			"Resource": "arn:aws:s3:::config-bucket-094724549126"
		},
		{
			"Sid": "AWSConfigBucketDelivery",
			"Effect": "Allow",
			"Principal": {
				"Service": "config.amazonaws.com"
			},
			"Action": "s3:PutObject",
			"Resource": "arn:aws:s3:::config-bucket-094724549126/AWSLogs/094724549126/Config/*",
			"Condition": {
				"StringEquals": {
					"s3:x-amz-acl": "bucket-owner-full-control"
				}
			}
		}
	]
}
POLICY
}


resource "aws_s3_bucket" "infralight-ido-all" {
	bucket = "infralight-ido-all"
	force_destroy = false
}