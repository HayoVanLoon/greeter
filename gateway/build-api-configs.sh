#!/usr/bin/env bash

set -eo pipefail

# Initialise with environment variables if present.
export PROJECT="${PROJECT}"
export PROJ_HASH="${PROJ_HASH}"
export GATEWAY_SVC_NAME
export BACKEND_SVC_1_NAME

usage() {
	echo "Well ain't that cute. BUT IT'S WRONG!
$1
"
	exit 1
}

process() {
	rm -f "${1}"
	cp "../${1}" . || exit 1
	sed -i -E "s/\\$\\{PROJECT\\}/${2}/g" "${1}"
	sed -i -E "s/\\$\\{PROJ_HASH\\}/${3}/g" "${1}"
	sed -i -E "s/\\$\\{GATEWAY_SVC_NAME\\}/${4}/g" "${1}"
	sed -i -E "s/\\$\\{BACKEND_SVC_1_NAME\\}/${5}/g" "${1}"
}

while getopts p:h:g:1: o; do
	echo "${o}" ${OPTARG}
	case "${o}" in
	"p") PROJECT=${OPTARG};;
	"h") PROJ_HASH=${OPTARG};;
	"g") GATEWAY_SVC_NAME=${OPTARG};;
	"1") BACKEND_SVC_1_NAME=${OPTARG};;
	*) usage ;;
	esac
done

echo "${PROJECT}" \
		"${PROJ_HASH}" \
		"${GATEWAY_SVC_NAME}" \
		"${BACKEND_SVC_1_NAME}"

if [ -z "${PROJECT}" ]; then
	usage "Missing '-p <project id>'"
fi
if [ -z "${PROJ_HASH}" ]; then
	usage "Missing '-h <cloud run project hash>'"
fi
if [ -z "${GATEWAY_SVC_NAME}" ]; then
	usage "Missing '-g <gateway service name>'"
fi
if [ -z "${BACKEND_SVC_1_NAME}" ]; then
	usage "Missing '-1 <backend service 1 name>'"
fi

OLD_DIR="${PWD}"
mkdir -p var
cd var || exit 1

for f in api_config.yaml api_config_auth.yaml api_config_http.yaml; do
	process "${f}" \
		"${PROJECT}" \
		"${PROJ_HASH}" \
		"${GATEWAY_SVC_NAME}" \
		"${BACKEND_SVC_1_NAME}"
done

cd "$OLD_DIR" || exit
