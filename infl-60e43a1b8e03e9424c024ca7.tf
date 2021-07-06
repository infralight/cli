resource "aws_ebs_volume" "vol-04be9b05d47f80fa5" {
	availability_zone = "us-east-2a"
	encrypted = false
	multi_attach_enabled = false
	size = 8
	snapshot_id = "snap-0c9864097499cc8a9"
	type = "gp2"
	iops = 100
}