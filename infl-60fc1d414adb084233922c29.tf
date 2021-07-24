resource "aws_kms_alias" "alias/cloud-trial-management-events" {
	name = "alias/cloud-trial-management-events"
	target_key_id = "a3b57fc9-4bcd-404e-af7e-4c815906d437"
}