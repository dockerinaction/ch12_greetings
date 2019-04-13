IMAGE_REPOSITORY := dockerinaction/ch12_greetings

build:
	docker build -t $(IMAGE_REPOSITORY):api -f api/Dockerfile .

push:
	docker push $(IMAGE_REPOSITORY):api

deploy-stack:
	docker stack deploy --compose-file docker-compose.yml greetings

SLEEP_TIME_IN_SECS := 10
destroy-stack:
	@echo "removing greetings stack"
	docker stack rm  greetings
	@echo "sleeping $(SLEEP_TIME_IN_SECS) seconds so Swarm can converge its state"
	sleep $(SLEEP_TIME_IN_SECS)
	@echo "removing orphaned greetings containers"
	docker container ls -a | grep greetings | awk '{ print $NF }' | xargs docker container rm -f
	@echo "sleeping $(SLEEP_TIME_IN_SECS) seconds so Swarm can converge its state"
