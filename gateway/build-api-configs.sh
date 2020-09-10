#!/usr/bin/env bash

set -eo pipefail

# Initialise with environment variables if present.
export PROJECT="${PROJECT}"
export PROJ_HASH="${PROJ_HASH}"
export CLIENT_ID="${CLIENT_ID}"
export GATEWAY_SVC_NAME
export BACKEND_SVC_1_NAME

usage() {
	echo "Well ain't that cute. BUT IT'S WRONG!

Usage:
	$(basename "${0}") [ARGS...]

Args:
	-p: project id (optional if env var PROJECT is set)
	-h: cloud run project hash (optional if env var PROJ_HASH is set)
	-g: gateway service name
	-1: backend service 1 name

$1"
	exit 1
}

process() {
	rm -f "${1}"
	cp "../${1}" . || exit 1
	sed -i -E "s/\\$\\{PROJECT\\}/${2}/g" "${1}"
	sed -i -E "s/\\$\\{PROJ_HASH\\}/${3}/g" "${1}"
	sed -i -E "s/\\$\\{CLIENT_ID\\}/${4}/g" "${1}"
	sed -i -E "s/\\$\\{GATEWAY_SVC_NAME\\}/${5}/g" "${1}"
	sed -i -E "s/\\$\\{BACKEND_SVC_1_NAME\\}/${6}/g" "${1}"
	sed -i -E "s/\\$\\{BACKEND_SVC_2_NAME\\}/${7}/g" "${1}"
}

while getopts p:h:c:g:1:2: o; do
	case "${o}" in
	"p") PROJECT=${OPTARG};;
	"h") PROJ_HASH=${OPTARG};;
	"c") CLIENT_ID=${OPTARG};;
	"g") GATEWAY_SVC_NAME=${OPTARG};;
	"1") BACKEND_SVC_1_NAME=${OPTARG};;
	"2") BACKEND_SVC_2_NAME=${OPTARG};;
	*) usage ;;
	esac
done

if [ -z "${PROJECT}" ]; then
	usage "Missing '-p <project id>'"
fi
if [ -z "${PROJ_HASH}" ]; then
	usage "Missing '-h <cloud run project hash>'"
fi
if [ -z "${CLIENT_ID}" ]; then
	usage "Missing '-c <oauth client id>'"
fi
if [ -z "${GATEWAY_SVC_NAME}" ]; then
	usage "Missing '-g <gateway service name>'"
fi
if [ -z "${BACKEND_SVC_1_NAME}" ]; then
	usage "Missing '-1 <backend service 1 name>'"
fi
if [ -z "${BACKEND_SVC_2_NAME}" ]; then
	usage "Missing '-2 <backend service 1 name>'"
fi

OLD_DIR="${PWD}"
mkdir -p var
cd var || exit 1

for f in api_config.yaml api_config_auth.yaml api_config_http.yaml; do
	process "${f}" \
		"${PROJECT}" \
		"${PROJ_HASH}" \
		"${CLIENT_ID}" \
		"${GATEWAY_SVC_NAME}" \
		"${BACKEND_SVC_1_NAME}" \
		"${BACKEND_SVC_2_NAME}"
done

cd "$OLD_DIR" || exit
