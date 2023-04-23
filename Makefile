
BuildFlag+= -X todo.codeVersion=$(BRANCH_NAME)_$(SHORT_SHA)
export LDFLAGS=$(BuildFlag)


proto:
	protoc -I=./proto --go-new_out=paths=source_relative:./proto ./proto/*.proto

.PHONY: proto