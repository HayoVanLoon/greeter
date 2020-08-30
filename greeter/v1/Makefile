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
MODULE_NAME := greeter
VERSION := $(shell basename `pwd`)
SERVICE_NAME := $(MODULE_NAME)-$(VERSION)
SERVICE_ACCOUNT := $(SERVICE_NAME)@$(PROJECT).iam.gserviceaccount.com

# GCP Region
REGION=europe-west1

# GCP Labels added to Cloud Run deployment (for tracking costs)
LABELS := app=cloud-endpoints-grpc

# Docker-related
IMAGE_NAME := $(MODULE_NAME)_$(VERSION)
TAG := latest

# Settings for local deployments
DEFAULT_PORT=8080


.PHONY: clean

check:
ifndef PROJECT
	$(error PROJECT not set)
endif

all: build push-gcr service-account deploy-unsafe

# TODO(hvl): unsafe...
release: build push-gcr deploy-unsafe

clean:
	go clean

run:
	export PORT=$(DEFAULT_PORT) && \
	go run server.go

build:
	docker build -t $(IMAGE_NAME) .

docker-run:
	docker run \
		--network="host" \
		$(IMAGE_NAME)

push-gcr: check
	docker tag $(IMAGE_NAME) gcr.io/$(PROJECT)/$(IMAGE_NAME):$(TAG)
	docker push gcr.io/$(PROJECT)/$(IMAGE_NAME)

service-account: check
	gcloud iam service-accounts create $(SERVICE_NAME) \
		--project=$(PROJECT)

deploy: check
	gcloud run deploy $(SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed \
		--memory=128Mi \
		--image=gcr.io/$(PROJECT)/$(IMAGE_NAME) \
		--no-allow-unauthenticated \
		--service-account=$(SERVICE_ACCOUNT) \
		--labels=$(LABELS)

deploy-unsafe: check
	gcloud run deploy $(SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed \
		--memory=128Mi \
		--image=gcr.io/$(PROJECT)/$(IMAGE_NAME) \
		--allow-unauthenticated \
		--service-account=$(SERVICE_ACCOUNT) \
		--labels=$(LABELS)

smoke-test-local:
	go run client/client.go -host=localhost -port=$(DEFAULT_PORT) -name=Jonas

smoke-test-cloud: check
	URL=$$(gcloud run services describe $(SERVICE_NAME) \
 		--project=$(PROJECT) \
 		--region=$(REGION) \
 		--platform=managed \
 		--format='value(status.address.url)') && \
	go run client/client.go -host="$${URL#https://}" -port=443 -skip-auth=true

destroy: check
	gcloud iam service-accounts delete $(SERVICE_ACCOUNT) \
		--project=$(PROJECT)
	gcloud run services delete $(SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed
