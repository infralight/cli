resource "aws_security_group_rule" "sgr-0012dfda2a4411c08" {
	type = "egress"
	from_port = 8200
	to_port = 8200
	protocol = "tcp"
	cidr_blocks = ["0.0.0.0/0"] 
	security_group_id = "sg-08742581b90800426"
	description = "Vault"
}