resource "aws_iam_policy" "ANPARMDQF5IDA2OFE6MHW" {
	name = "AllowAllSqsInDevelopmentRegion"
	path = "/"
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "VisualEditor0",
			"Effect": "Allow",
			"Action": "sqs:*",
			"Resource": "*",
			"Condition": {
				"StringEquals": {
					"aws:RequestedRegion": "us-west-2"
				}
			}
		}
	]
}
POLICY
}


resource "aws_iam_policy" "ANPARMDQF5IDDJ2E3KQH6" {
	name = "AllowCRUDAndAssumeRole"
	path = "/"
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "VisualEditor0",
			"Effect": "Allow",
			"Action": [
				"sts:AssumeRole",
				"iam:DetachRolePolicy",
				"iam:CreateRole",
				"iam:AttachRolePolicy"
			],
			"Resource": "*"
		}
	]
}
POLICY
}