resource "aws_iam_role" "AROARMDQF5IDGRPWA3JNJ" {
	unique_id = "AROARMDQF5IDGRPWA3JNJ"
	name = "stag-dragonfly-cluster-C30TnJdV20210909135738747400000002"
	assume_role_policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "EKSWorkerAssumeRole",
			"Effect": "Allow",
			"Principal": {
				"Service": "ec2.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}
	]
}
POLICY
	max_session_duration = 3600
	tags = {
		"Environment" = "stag"
	}
}