.PHONY: build

build:
	go mod tidy
	go mod vendor
	go build -o capten ./cmd/main.go 

.PHONY: build.linux
build.linux:
	@echo "👷👷 Building Capten CLI"
	@rm -rf capten
	@mkdir capten
	@go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o capten/capten cmd/main.go

.PHONY: build.mac
build.mac:
	@echo "👷👷 Building Capten CLI"
	@rm -rf capten
	@mkdir capten
	@go mod download && GOOS=darwin GOARCH=amd64 go build -o capten/capten.app cmd/main.go

.PHONY: build.release-linux
build.release-linux: build.linux
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

	# Download and extract Terraform binary for linux
	@curl -LO https://releases.hashicorp.com/terraform/0.12.31/terraform_0.12.31_linux_amd64.zip
	@unzip terraform_0.12.31_linux_amd64.zip -d capten/ && mv capten/terraform capten/terraform-linux
	@chmod +x capten/terraform-linux
	@rm terraform_0.12.31_linux_amd64.zip

	# Download and extract Talosctl binary for linux
	@curl -LO https://github.com/siderolabs/talos/releases/download/v1.4.8/talosctl-linux-amd64
	@mv talosctl-linux-amd64 capten/terraform_modules/talosctl
	@chmod +x capten/terraform_modules/talosctl

	@zip -r capten.zip capten/*
	# remove this release folder as ci pipeline is complaining
	@rm -rf capten
	@echo "✅ Release Build Complete ✅"

.PHONY: build.release-mac
build.release-mac: build.mac
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

	# Download and extract Terraform binary for mac
	@curl -LO https://releases.hashicorp.com/terraform/0.12.31/terraform_0.12.31_darwin_amd64.zip
	@unzip terraform_0.12.31_darwin_amd64.zip -d capten/ && mv capten/terraform capten/terraform-mac
	@chmod +x capten/terraform-mac
	@rm terraform_0.12.31_darwin_amd64.zip

	# Download and extract Talosctl binary for mac
	@curl -LO https://github.com/siderolabs/talos/releases/download/v1.4.8/talosctl-darwin-amd64
	@mv talosctl-darwin-amd64 capten/terraform_modules/talosctl
	@chmod +x capten/terraform_modules/talosctl

	@zip -r capten.zip capten/*
	# remove this release folder as ci pipeline is complaining
	@rm -rf capten
	@echo "✅ Release Build Complete ✅"