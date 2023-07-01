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

.PHONY: build.release
build.release: build.all
	@echo "👷👷 Building release ..."
	@mkdir release
	@mkdir release/config release/cert release/terraform_modules release/apps release/templates

	# move binaries and configs
	@mv capten-* release
	@cp -rf config/* release/config/

	# move templates
	@cp -rf templates/* release/templates/

	# git pull dataplane repo, versioning needs to be confirmed
	@git clone https://github.com/kube-tarian/controlplane-dataplane.git
	@cp -rf controlplane-dataplane/* release/terraform_modules/
	@rm -rf controlplane-dataplane

	# move apps
	@mkdir release/apps/values
	@cp -rf apps/* release/apps/

	# copy readme
	@cp README.md release/README.md

	@zip -r capten.zip release
	@echo "✅ Release Build Complete ✅"