resource "aws_kms_key" "Default key that protects my Lightsail signing keys when no other key is defined" {
	description = "Default key that protects my Lightsail signing keys when no other key is defined"
	key_usage = "ENCRYPT_DECRYPT"
	customer_master_key_spec = "SYMMETRIC_DEFAULT"
	is_enabled = true
	enable_key_rotation = false
}