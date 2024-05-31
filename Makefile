run:
	docker build -t plain.do . && docker run -d -v .:/plain.do -p 8080:8080 plain.do

run-dev:
	docker build --build-arg GISN_ENV=development -t plain.do . && docker run -v .:/plain.do plain.do

stop:
	docker ps -q --filter "ancestor=plain.do" | xargs docker stop

shell:
	docker exec -it $$(docker ps -q --filter "ancestor=plain.do" | head -n 1) /bin/bash

logs:
	docker logs -f $$(docker ps -q --filter "ancestor=plain.do" | head -n 1)

migrate:
	docker exec $$(docker ps -q --filter "ancestor=plain.do" | head -n 1) go run utils/migrate.go
