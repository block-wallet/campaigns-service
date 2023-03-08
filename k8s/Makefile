

ecr/create:
	@rm -rf .terraform*
	@terraform init \
	-backend-config="profile=blank-devops" \
	-backend-config="region=us-east-2" \
	-backend-config="bucket=blank-terraform-states" \
	-backend-config="key=projects/ecr/go-template.tfstate"
	@terraform apply 

ecr/destroy:
	@rm -rf .terraform*
	@terraform init \
	-backend-config="profile=blank-devops" \
	-backend-config="region=us-east-2" \
	-backend-config="bucket=blank-terraform-states" \
	-backend-config="key=projects/ecr/go-template.tfstate"
	@terraform destroy