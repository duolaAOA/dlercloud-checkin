.PHONY: build clean

GO 				:= GO111MODULE=on go
GO_FILES 		:= $(shell git ls-files '*.go' | grep -v '^vendor/')

linux:
		@echo "building..."
		GOOS=linux GOARCH=amd64 ${GO} build -ldflags "${GO_LDFLAGS}" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" -o ./bin/dler ./*.go
		$(if $(shell command -v upx), upx ./bin/dler)

darwin:
		@echo "building..."
		GOOS=darwin GOARCH=amd64 ${GO} build -ldflags "${GO_LDFLAGS}" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" -o ./bin/dler ./*.go
		$(if $(shell command -v upx), upx ./bin/dler)

clean:
		@echo "cleaning..."
		rm -rf bin/
		rm -rf build/