resource "aws_cloudfront_distribution" "E2IYEU59P9HSO4" {
	aliases = ["stag.infralight.cloud"] 
	default_root_object = "index.html"
	is_ipv6_enabled = true
	http_version = "http2"
	enabled = true
	price_class = "PriceClass_All"
	restrictions {
		geo_restriction {
			locations = ["IR"] 
			restriction_type = "blacklist"
		}
	}
	viewer_certificate {
		acm_certificate_arn = "arn:aws:acm:us-east-1:094724549126:certificate/9081ac57-8983-4a61-8a5e-21d60cc8d79f"
		cloudfront_default_certificate = false
		iam_certificate_id = ""
		minimum_protocol_version = "TLSv1.2_2019"
		ssl_support_method = "sni-only"
	}
	web_acl_id = "arn:aws:wafv2:us-east-1:094724549126:global/webacl/cloudfront-waf/e2228bd9-f65b-47c4-8b81-ff49e3276c68"
	custom_error_response {
		error_caching_min_ttl = 300
		error_code = 400
		response_code = "200"
		response_page_path = "/index.html"
	}
	custom_error_response {
		error_caching_min_ttl = 300
		error_code = 403
		response_code = "200"
		response_page_path = "/index.html"
	}
	logging_config {
		bucket = "aws-cloudtrail-logs-094724549126-0b9a848c.s3.amazonaws.com"
		include_cookies = false
		prefix = "stag-cf-logs/"
	}
	origin {
		domain_name = "stag-infralight-dashboard-static.s3.amazonaws.com"
		origin_id = "S3-stag-infralight-dashboard-static"
		origin_path = ""
		s3_origin_config {
			origin_access_identity = "origin-access-identity/cloudfront/E2VFLX4T5MO11V"
		}
		custom_origin_config {
			http_port = 0
			https_port = 0
			origin_protocol_policy = ""
			origin_ssl_protocols = [] 
			origin_keepalive_timeout = 0
			origin_read_timeout = 0
		}
	}
	default_cache_behavior {
		allowed_methods = ["HEAD","DELETE","POST","GET","OPTIONS","PUT","PATCH"] 
		cached_methods = ["HEAD","GET"] 
		cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6"
		compress = false
		default_ttl = 0
		field_level_encryption_id = ""
		forwarded_values {
			query_string = false
			query_string_cache_keys = [] 
			headers = []
			cookies {
				forward = "none"
				whitelisted_names = [] 
			}
		}
		max_ttl = 0
		min_ttl = 0
		origin_request_policy_id = ""
		realtime_log_config_arn = ""
		smooth_streaming = false
		target_origin_id = "S3-stag-infralight-dashboard-static"
		viewer_protocol_policy = "redirect-to-https"
	}
}