# Where proto files live
PROTO_DIR := ../api/protos

# Where Go code should be generated (usermanagement service)
GO_OUT := ../api/protos

# Where the Envoy descriptor should be generated
DESC_OUT := ../envoy/protos/api.pb

# Google APIs for annotations, etc.
GOOGLEAPIS_DIR := ../usermanagement/third_party/googleapis

# All proto files
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

all: protos

protos:
	@echo "Ensuring Google APIs are present..."
	@if [ ! -d "$(GOOGLEAPIS_DIR)" ]; then \
		echo "Cloning googleapis repository..."; \
		git clone https://github.com/googleapis/googleapis.git $(GOOGLEAPIS_DIR); \
	fi

	@echo "Generating Go code + Envoy descriptor..."
	protoc \
	  -I $(PROTO_DIR) \
	  -I $(GOOGLEAPIS_DIR) \
	  --go_out=$(GO_OUT) --go_opt=paths=source_relative \
	  --go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
	  --descriptor_set_out=$(DESC_OUT) \
	  --include_imports \
	  --include_source_info \
	  $(PROTO_FILES)

	@echo "Protos and descriptor generated"

clean:
	@echo "Cleaning generated Go files and descriptor..."
	find $(GO_OUT) -name "*.pb.go" -delete
	rm -f $(DESC_OUT)