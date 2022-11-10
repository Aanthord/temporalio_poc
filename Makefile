all: kafka telemetry coston airbyte watson

kafka:
	docker-compose -f ./compose-files/kafka-stack-docker-compose/full-stack.yml

telemetry:
	for X in jaeger prometheus
	do
		docker-compose -f ./compose-files/$X/docker-compose.yml up && sleep 10 && echo $X is complete
	done

coston:
	docker-compose -f ./compose-files/coston/docker-compose.yml


airbyte:
	docker-compose -f ./compose-files/airbyte/docker-compose.yml

watson:
	docker-compose -f ./compose-files/watson/docker-compose.yml

run:

connect:
