resource "aws_instance" "i-029c0a740a1a0a6f1" {
	ami = "ami-0ebe1b93841528ab3"
	associate_public_ip_address = true
	availability_zone = "us-west-2a"
	cpu_core_count = 2
	cpu_threads_per_core = 1
	ebs_optimized = false
	hibernation = false
	iam_instance_profile = "stag-vault"
	instance_type = "t2.large"
	key_name = "stag"
	private_ip = "10.1.0.198"
	source_dest_check = true
	subnet_id = "subnet-093a001e537625da2"
	vpc_security_group_ids = ["sg-0aee2f15547d2f2ce"] 
	metadata_options {
		http_endpoint = "enabled"
		http_put_response_hop_limit = 1
		http_tokens = "optional"
	}
}