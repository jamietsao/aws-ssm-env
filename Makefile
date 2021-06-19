# credit: https://vic.demuzere.be/articles/golang-makefile-crosscompile/
PLATFORMS := darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64

checkenv:
ifndef TAG
	$(error release TAG is required - e.g v0.1.0)
endif

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
name = aws-ssm-env
longname = $(name)-$(TAG)-$(os)-$(arch)

release: $(PLATFORMS)

$(PLATFORMS): checkenv
	GOOS=$(os) GOARCH=$(arch) go build -o 'bin/$(longname)/$(name)' main.go
	cd bin/$(longname) && zip $(longname).zip $(name)

.PHONY: checkenv release $(PLATFORMS) zip
