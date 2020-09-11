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

# Python Environment Activation
V := . .venv/bin/activate
MAKE := make

# Naming
MODULE_NAME := bartender
VERSION := $(shell basename `pwd`)
SERVICE_NAME := $(MODULE_NAME)-$(VERSION)
SERVICE_ACCOUNT := $(SERVICE_NAME)@$(PROJECT).iam.gserviceaccount.com

# Protocol Buffers
PROTO_APIS_DIR := ../../apis

# GCP Region
REGION=europe-west1

# GCP Labels added to Cloud Run deployment (for tracking costs)
LABELS := app=cloud-endpoints-grpc

# Docker-related
IMAGE_NAME := $(MODULE_NAME)_$(VERSION)
TAG := latest

# Settings for local deployments
DEFAULT_PORT=8080

GATEWAY := https://esp-v2beta-$(PROJ_HASH)-ew.a.run.app


.PHONY: clean

check:
ifndef PROJECT
	$(error PROJECT not set)
endif

all: build push-gcr service-account deploy
	echo "Manual step: Set up OAuth2 Client ID at: https://console.cloud.google.com/apis/credentials"

release: build push-gcr deploy

protoc:
	$(MAKE) -C $(PROTO_APIS_DIR) protoc-py

venv: protoc
	rm -rf .venv
	virtualenv -p python3 .venv
	$(V); pip install $(PROTO_APIS_DIR)/python
	$(V); pip install -r requirements.txt

freeze:
	$(V); pip freeze

clean:
	rm -rf var
	rm -rf .pytest_cache

test: clean
	$(V); pytest

run:
	export PORT=$(DEFAULT_PORT) && \
	$(V); python src/server.py

# Forego complexity of building multiple images; copy everything over
build: protoc
	rm -rf var
	mkdir -p var/var
	cp -r $(PROTO_APIS_DIR)/python/var/hayovanloon var/var
	cp $(PROTO_APIS_DIR)/python/setup.py var
	docker build -t $(IMAGE_NAME) .

docker-run:
	docker run \
		--network="host" \
		--env PORT=$(PORT) \
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

do-create-beer-local:
	grpcurl -plaintext \
		-d '{"brand": "BrewDog", "beer":{"name": "ElvisJuice", "bottle_cl": 33}}' \
		localhost:8080 \
		hayovanloon.bartender.v1.Bartender/CreateBeer

do-get-beer-local:
	grpcurl -plaintext \
		-d '{"name": "ElvisJuice", "brand": "BrewDog"}' \
		localhost:8080 \
		hayovanloon.bartender.v1.Bartender/GetBeer

do-create-beer:
	TOKEN=$$(gcloud auth print-identity-token \
		--audiences=$(CLIENT_ID) --impersonate-service-account=$(DEV_SA)) && \
	curl \
		-d '{"name": "ElvisJuice", "bottle_cl": 33}' \
		-H "Authorization: Bearer $$TOKEN" \
		$(GATEWAY)/v1/beers/BrewDog/type

do-list-beers:
	TOKEN=$$(gcloud auth print-identity-token \
		--audiences=$(CLIENT_ID) --impersonate-service-account=$(DEV_SA)) && \
	curl \
		-H "Authorization: Bearer $$TOKEN" \
		$(GATEWAY)/v1/beers

smoke-test-cloud: check
	URL=$$(gcloud run services describe $(SERVICE_NAME) \
 		--project=$(PROJECT) \
 		--region=$(REGION) \
 		--platform=managed \
 		--format='value(status.address.url)') && \
	echo "TODO: add grpcurl call to $${URL#https://}"

destroy: check
	gcloud iam service-accounts delete $(SERVICE_ACCOUNT) \
		--project=$(PROJECT)
	gcloud run services delete $(SERVICE_NAME) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed
