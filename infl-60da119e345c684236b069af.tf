resource "aws_iam_policy" "ANPARMDQF5IDMBTEAW6A6" {
	name = "AllowAllLambdaInDevelopmentRegion"
	path = "/"
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "VisualEditor0",
			"Effect": "Allow",
			"Action": "lambda:*",
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