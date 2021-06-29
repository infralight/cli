resource "aws_s3_bucket" "nops-report-bucket" {
	bucket = "nops-report-bucket"
	force_destroy = false
	policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"AWS": "arn:aws:iam::386209384616:root"
			},
			"Action": [
				"s3:GetBucketAcl",
				"s3:GetBucketPolicy"
			],
			"Resource": "arn:aws:s3:::nops-report-bucket"
		},
		{
			"Effect": "Allow",
			"Principal": {
				"AWS": "arn:aws:iam::386209384616:root"
			},
			"Action": "s3:PutObject",
			"Resource": "arn:aws:s3:::nops-report-bucket/*"
		}
	]
}
POLICY
}