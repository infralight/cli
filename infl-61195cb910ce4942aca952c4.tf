resource "aws_iam_role" "AROARMDQF5IDFH7QAVLCM" {
	unique_id = "AROARMDQF5IDFH7QAVLCM"
	name = "AmazonSSMRoleForAutomationAssumeQuickSetup"
	assume_role_policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Service": "ssm.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}
	]
}
POLICY
	max_session_duration = 3600
}