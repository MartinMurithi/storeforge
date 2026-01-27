# Makefile - Usermanagement proto generation

# Directories
PROTO_DIR := proto
THIRD_PARTY := third_party/googleapis

# Plugins (installed via `go install`)
PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc)

# Proto files
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

# Default target
all: protos

# Compile all protos
protos:
	@echo "Generating Go code from proto files..."
ifeq ($(PROTOC_GEN_GO),)
	$(error "protoc-gen-go not found! Run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest")
endif
ifeq ($(PROTOC_GEN_GO_GRPC),)
	$(error "protoc-gen-go-grpc not found! Run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
endif
	protoc \
	  -I $(PROTO_DIR) \
	  -I $(THIRD_PARTY) \
	  --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
	  --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
	  $(PROTO_FILES)
	@echo "Proto compilation complete!"

# Clean generated code
clean:
	@echo "Removing generated .pb.go files..."
	find $(PROTO_DIR) -name "*.pb.go" -delete
	@echo "Clean complete"
