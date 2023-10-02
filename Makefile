.PHONY: build

build:
	go mod tidy
	go mod vendor
	go build -o capten ./cmd/main.go 

.PHONY: build.all
build.all:
	@echo "👷👷 Building Capten CLI"
	@rm -rf capten
	@mkdir capten
	@go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o capten/capten cmd/main.go
#	@go mod download && GOOS=darwin GOARCH=amd64 go build -o capten/capten.app cmd/main.go
#	@go mod download && GOOS=windows GOARCH=amd64 go build -o capten/capten.exe cmd/main.go

.PHONY: build.release
build.release: build.all
	@echo "👷👷 Building release ..."
	@mkdir capten/config capten/cert capten/terraform_modules capten/apps capten/templates

	# move configs
	@cp -rf config/* capten/config/

	# move templates
	@cp -rf templates/* capten/templates/

	# git pull dataplane repo, versioning needs to be confirmed
	@git clone https://github.com/kube-tarian/controlplane-dataplane.git
	@cp -rf controlplane-dataplane/* capten/terraform_modules/
	@rm -rf controlplane-dataplane

	# move apps
	@mkdir capten/apps/values
	@cp -rf apps/* capten/apps/

	# copy readme
	@cp README.md capten/README.md

	# make all scripts executable
	@find ./capten/ -type f -name "*.sh" -exec chmod +x {} \;

	# Download and extract Terraform binary
	@curl -LO https://releases.hashicorp.com/terraform/0.12.31/terraform_0.12.31_linux_amd64.zip
	@unzip terraform_0.12.31_linux_amd64.zip -d capten/
	@chmod +x capten/terraform
	@rm terraform_0.12.31_linux_amd64.zip

# Download talosctl binary
	@curl -LO https://github.com/siderolabs/talos/releases/download/v1.5.2/talosctl-linux-amd64
	@mv talosctl-linux-amd64 capten/config
	@chmod +x capten/terraform_modules/azure/talos/talosctl

	@zip -r capten.zip capten/*
	# remove this release folder as ci pipeline is complaining
	@rm -rf capten
	@echo "✅ Release Build Complete ✅"