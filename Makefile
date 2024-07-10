run:
	docker build -t plain.do . && docker run -d -v .:/plain.do -p 8080:8080 plain.do

stop:
	docker ps -q --filter "ancestor=plain.do" | xargs docker stop

shell:
	docker exec -it $$(docker ps -q --filter "ancestor=plain.do" | head -n 1) /bin/bash

logs:
	docker logs -f $$(docker ps -q --filter "ancestor=plain.do" | head -n 1)
