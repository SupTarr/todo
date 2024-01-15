build:
	go build \
		-ldflags "-X ,sim.buildcommit=`git rev-parse --short HEAD` \
		-X main.buildtime=`date "+%Y-%m-%dT%H:%S%Z:00"`" \
		-o todos

run:
	./todos

postgres:
	docker run --name todos -e POSTGRES_PASSWORD=2009 -e POSTGRES_USER=postgres -e POSTGRES_DB=todos -d -p 5432:5432 postgres

reload:
	echo "GET http://:808/limitz" | vegeta attack -rate=10/s -duration=1s | vegeta report
