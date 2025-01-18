TAG := latest
DOCKER_IMAGE := go-license-management:${TAG}

build:
	docker build -t ${DOCKER_IMAGE} .


push-local:
	docker push localhost:5000/${DOCKER_IMAGE}


push-image:
	docker push ${DOCKER_IMAGE}