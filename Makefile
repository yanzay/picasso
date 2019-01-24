DATE=`date +'%Y%m%d%H%M%S'`

build:
	go build -v

run: build
	./picasso --local

install:
	go install -v

test:
	@rm -f _test.db
	@go test -v ./...

backup:
	@mkdir -p backups
	@cp data.db backups/data.$(DATE).bak
	@echo backups/data.$(DATE).bak created

generate: backup
	go run cmd/playersgen/main.go

guilds: backup
	go run cmd/guilds/main.go

lint:
	gometalinter --enable-all ./...

clean:
	rm data.db
	rm picasso
