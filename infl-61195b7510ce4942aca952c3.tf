resource "aws_ebs_volume" "vol-09efba780c5a176e2" {
	availability_zone = "us-west-2a"
	encrypted = false
	iops = 100
	multi_attach_enabled = false
	size = 8
	snapshot_id = "snap-0875a8226ee4e843b"
	type = "gp2"
}