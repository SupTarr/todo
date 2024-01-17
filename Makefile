build:
	go build \
		-ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
		-X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
		-o todo

run:
	./todo

maria:
	docker run -p localhost:3306:3306 --name some-mariadb \
	-e MARIADB_ROOT_PASSWORD=my-secret-pw -e MARIADB_DATABASE=myapp -d mariadb:latest

reload:
	echo "GET http://localhost:8081/limitz" | vegeta attack -rate=10/s -duration=1s | vegeta report

image:
	docker build -t todo:test -f Dockerfile .

container:
	docker run -p:8081:8081 --env-file ./local.env \
	--link some-mariadb:db --name runningtodo todo:test
