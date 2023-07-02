.PHONY: build

build:
	go mod tidy
	go mod vendor
	go build -o capten ./cmd/main.go 

.PHONY: build.all
build.all:
	@echo "👷👷 Building captain clis"
	@go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o capten-linux cmd/main.go
	@go mod download && GOOS=darwin GOARCH=amd64 go build -o capten-macos cmd/main.go
	@go mod download && GOOS=windows GOARCH=amd64 go build -o capten-windows.exe cmd/main.go

.PHONY: certs.generate
certs.generate:
	@echo "👷👷 Generating certs"
	@go build -o cert-gen pkg/cert/cmd/main.go
	@./cert-gen
	@rm cert-gen

.PHONY: build.release
build.release: build.all certs.generate
	@echo "👷👷 Building release ..."
	@mkdir release
	@mkdir release/config release/certs release/terraform release/apps release/templates

	# move binaries and configs
	@mv capten-* release
	@cp config/aws_config.yaml release/config/aws_config.yaml

	# move templates
	@cp -rf templates/* release/templates/

	# git pull dataplane repo, versioning needs to be confirmed
	@git clone https://github.com/kube-tarian/controlplane-dataplane.git
	@mkdir release/terraform/aws
	@cp -rf controlplane-dataplane/aws/* release/terraform/aws/
	@rm -rf controlplane-dataplane

	# move apps
	@mkdir release/apps/values
	@cp config/apps.yaml release/apps/apps.yaml
	@cp -rf config/values/* release/apps/values

	# generate and move certs
	@mv cert/* release/certs

	# copy readme
	@cp README.md release/README.md

	@zip -r capten.zip release
	@echo "✅ Release Build Complete ✅"