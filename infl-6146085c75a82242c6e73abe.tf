resource "aws_lambda_function" "lambda-warmer" {
	function_name = "lambda-warmer"
	role = "arn:aws:iam::094724549126:role/prod-lambda-warmer"
	handler = "lambda_function.lambda_handler"
	source_code_hash = "Lg8cSomJAnnn8JjsDxY8wLKu+ogAcRDOWsizpqbD/vU="
	runtime = "python3.6"
	timeout = 15
	memory_size = 128
	environment {
		variables = {
			INVOKED_LAMBDA = "prod-reverse-learning"
		}
	}
	tracing_config {
		mode = "PassThrough"
	}
}