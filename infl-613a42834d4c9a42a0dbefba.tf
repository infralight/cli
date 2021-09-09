resource "aws_iam_policy" "ANPARMDQF5IDJWAKUPJVK" {
	name = "eks-fargate-logging-policy"
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"logs:CreateLogStream",
				"logs:CreateLogGroup",
				"logs:DescribeLogStreams",
				"logs:PutLogEvents"
			],
			"Resource": "*"
		}
	]
}
POLICY
}


resource "aws_iam_policy" "ANPARMDQF5IDFXVYKXX6L" {
	name = "stag-dragonfly-cluster-C30TnJdV-elb-sl-role-creation20210909133750512400000001"
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "",
			"Effect": "Allow",
			"Action": [
				"ec2:DescribeInternetGateways",
				"ec2:DescribeAddresses",
				"ec2:DescribeAccountAttributes"
			],
			"Resource": "*"
		}
	]
}
POLICY
}