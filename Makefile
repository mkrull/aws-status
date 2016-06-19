.PHONY = clean builder runner

GOPROJ = github.com/mkrull/aws-status/sloppy-aws
MOUNTDIR = github.com/mkrull/aws-status

runner: build
	docker build -t sloppy-aws-runner -f Dockerfile $(PWD)

build: $(PWD)/build/sloppy-aws

$(PWD)/build:
	@mkdir build

$(PWD)/build/sloppy-aws: $(PWD)/build builder
	@docker run \
		-e GOPATH=/go \
		-v $(PWD):/go/src/$(MOUNTDIR) \
		-v $(PWD)/build:/go/bin \
		aws-check-builder sh -c "go install $(GOPROJ)"

builder:
	@docker build -t sloppy-aws-builder -f Dockerfile.builder $(PWD)

test: builder
	@docker run \
		-e GOPATH=/go \
		-v $(PWD):/go/src/$(MOUNTDIR) \
		aws-check-builder sh -c "cd /go/src/$(GOPROJ) && go test $$(go list ./... | grep -v /vendor/)"

clean:
	@rm -rf $(PWD)/build

run:
	@docker run sloppy-aws-runner
