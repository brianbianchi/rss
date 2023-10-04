build:
	go build -o bin/web web/*.go
	go build -o bin/cron cron/main.go
	go build -o bin/init db/init.go	

serve:
	go build -o bin/web web/*.go
	./bin/web

cj:
	go build -o bin/cron cron/main.go
	./bin/cron

init:
	go build -o bin/init db/init.go	
	./bin/init

clean:
	rm -rf bin
