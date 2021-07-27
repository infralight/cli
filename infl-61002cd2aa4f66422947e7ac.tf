resource "aws_security_group" "eks-cluster-sg-stag-ArangoCluster-704653202" {
	name = "eks-cluster-sg-stag-ArangoCluster-704653202"
	description = "EKS created security group applied to ENI that is attached to EKS Control Plane master nodes, as well as any managed workloads."
	vpc_id = "vpc-06af93201036c0670"
	ingress {
		protocol = "-1"
		cidr_blocks = [] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = ["","",""] 
	}
	egress {
		from_port = 80
		to_port = 80
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	egress {
		from_port = 0
		to_port = 65535
		protocol = "tcp"
		cidr_blocks = ["10.1.0.0/16"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = ["","","",""] 
	}
	egress {
		from_port = 8200
		to_port = 8200
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	egress {
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	egress {
		from_port = 8071
		to_port = 8071
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	egress {
		from_port = 22
		to_port = 22
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	egress {
		from_port = 8529
		to_port = 8529
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	egress {
		from_port = 0
		to_port = 0
		protocol = "tcp"
		cidr_blocks = ["192.168.232.0/21"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	egress {
		from_port = 443
		to_port = 443
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"] 
		ipv6_cidr_blocks = [] 
		prefix_list_ids = [] 
		security_groups = [] 
	}
	owner_id = "094724549126"
	tags = {
		"Name" = "eks-cluster-sg-stag-ArangoCluster-704653202"
		"aws:eks:cluster-name" = "stag-ArangoCluster"
		"kubernetes.io/cluster/stag-ArangoCluster" = "owned"
	}
}