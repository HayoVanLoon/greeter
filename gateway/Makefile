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

# Naming
MODULE_NAME := esp
VERSION := v2beta
SERVICE_NAME := $(MODULE_NAME)-$(VERSION)
SERVICE_ACCOUNT := $(SERVICE_NAME)@$(PROJECT).iam.gserviceaccount.com

# Backend services
GREETER_SERVICE_NAME := greeter-v1
GREETER_SERVICE_PROTO := ../proto/hayovanloon/greeter/greeter.proto

# GCP Region
REGION=europe-west1

# GCP Labels added to Cloud Run deployment (for tracking costs)
LABELS := app=cloud-endpoints-grpc

# PROJ_HASH is the short hash after the service name in the Cloud Run URL
PROJ_HASH :=
SELF_HOST := $(SERVICE_NAME)-$(PROJ_HASH)-ew.a.run.app

# Extensible Service Proxy parameters
ESP_VERSION := 2.16.0

# Path to a clone of the GoogleApis repository
PROTO_GOOGLE_APIS :=


.PHONY: clean

release: build deploy

check:
ifndef PROJ_HASH
	$(error PROJ_HASH not set)
endif
ifndef PROJECT
	$(error PROJECT not set)
endif

protoc-desc:
ifndef PROTO_GOOGLE_APIS
	$(error Missing PROTO_GOOGLE_APIS)
endif
	mkdir -p var
	protoc \
		--include_imports \
		--include_source_info \
      	--descriptor_set_out var/api_descriptor.pb \
		-I$(PROTO_GOOGLE_APIS) \
    	-I../proto \
    	$(GREETER_SERVICE_PROTO)

deploy-endpoints: check protoc-desc
	gcloud endpoints services deploy \
		var/api_descriptor.pb \
		api_config.yaml \
		api_config_http.yaml \
		--project=$(PROJECT)

build: check deploy-endpoints
	@export SVC=$$(gcloud endpoints services list | grep $(SERVICE_NAME)); \
	export SVC_VERSION=$$(gcloud endpoints configs list --service="$$SVC" | head -n 2 | tail -n 1 | sed -E "s/\s.+//g"); \
	echo Endpoints version: $${SVC_VERSION}; \
	vendor/gcloud_build_image.sh -s $(SELF_HOST) -c $${SVC_VERSION} -p $(PROJECT)

service-account: check
	-gcloud iam service-accounts create $(SERVICE_NAME) \
		--project=$(PROJECT)
	gcloud run services add-iam-policy-binding $(GREETER_SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed \
		--member="serviceAccount:$(SERVICE_ACCOUNT)" \
		--role="roles/run.invoker"

# As far as I understand: only sets up the service as a placeholder.
dummy-deploy: check
	gcloud run deploy $(SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed \
		--image="gcr.io/endpoints-release/endpoints-runtime-serverless:2" \
		--service-account=$(SERVICE_ACCOUNT) \
		--allow-unauthenticated \
		--labels=$(LABELS)

deploy: check
	@export SVC=$$(gcloud endpoints services list | grep $(SERVICE_NAME)); \
	export SVC_VERSION=$$(gcloud endpoints configs list --service="$$SVC" | head -n 2 | tail -n 1 | sed -E "s/\s.+//g"); \
	echo Endpoints version: $${SVC_VERSION}; \
	gcloud run deploy $(SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed \
		--image="gcr.io/$(PROJECT)/endpoints-runtime-serverless:$(ESP_VERSION)-$(SELF_HOST)-$${SVC_VERSION}" \
		--service-account=$(SERVICE_ACCOUNT) \
		--allow-unauthenticated \
		--set-env-vars=ESPv2_ARGS=--enable_debug \
		--labels=$(LABELS)

destroy: check
	gcloud iam service-accounts delete $(SERVICE_ACCOUNT) \
		--project=$(PROJECT)
	gcloud run services delete $(SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed