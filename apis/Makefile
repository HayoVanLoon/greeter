# Copyright 2020 Hayo van Loon
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROTO_ROOT := ./proto
PROTO_FILES := $(shell find $(PROTO_ROOT) -name '*.proto' | sort)

# GENPROTO_DIR is the (multi-tenant) directory for generated Go code.
GO_OUT := $(GENPROTO_DIR)

PY_OUT := python/var


check:
ifndef PROTO_GOOGLE_APIS
	$(error Missing PROTO_GOOGLE_APIS)
endif
ifndef GO_OUT
	$(error Missing GO_OUT)
endif

protoc-go: check
	export GO111MODULE=on; \
	protoc \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		-I$(PROTO_GOOGLE_APIS) \
		-I$(PROTO_ROOT) \
		$(PROTO_FILES)

quick-push:
ifndef GENPROTO_DIR
	$(error Missing GENPROTO_DIR)
endif
	-cd $(GENPROTO_DIR); ./quick-push.sh

protoc-py: check
	mkdir -p $(PY_OUT)
	rm -rf $(PY_OUT)/*
	. python/.venv/bin/activate; \
	python -m grpc.tools.protoc \
		--python_out=$(PY_OUT) \
		--grpc_python_out=$(PY_OUT) \
		-I$(PROTO_GOOGLE_APIS) \
		-I$(PROTO_ROOT) \
		$(PROTO_FILES)
