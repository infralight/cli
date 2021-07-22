resource "aws_iam_policy" "ANPARMDQF5IDC4YOSCWMN" {
	name = "stag-consumer-tfc_state_updater"
	path = "/"
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": [
				"ec2:CreateNetworkInterface",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DeleteNetworkInterface",
				"sts:AssumeRole"
			],
			"Effect": "Allow",
			"Resource": "*"
		},
		{
			"Action": [
				"sqs:ReceiveMessage",
				"sqs:DeleteMessage",
				"sqs:Get*",
				"sqs:List*"
			],
			"Effect": "Allow",
			"Resource": "arn:aws:sqs:us-west-2:094724549126:stag-producer-tfc_state_updater"
		}
	]
}
POLICY
}