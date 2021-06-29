resource "aws_iam_policy" "ANPARMDQF5IDDLFN442G4" {
	name = "AllowAllCloudFormartionInDevelopmentRegion"
	path = "/"
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "VisualEditor0",
			"Effect": "Allow",
			"Action": "cloudformation:*",
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