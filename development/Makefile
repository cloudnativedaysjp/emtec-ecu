.PHONY: up down rebuild dk-mock-server.yaml

up: dk-mock-server.yaml
	docker compose up

up-without-server: dk-mock-server.yaml
	docker compose -f docker-compose.without-server.yml up

down:
	docker compose down

rebuild:
	docker compose build --no-cache

dk-mock-server.yaml: dk-mock-server-base.yaml
	sed -e "s/DATE/$(shell date "+%Y-%m-%d/")" \
		-e "s/TIME1/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "1 minute")+09:00"/" \
		-e "s/TIME2/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "3 minute")+09:00"/" \
		-e "s/TIME3/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "5 minute")+09:00"/" \
		-e "s/TIME4/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "7 minute")+09:00"/" \
		dk-mock-server-base.yaml > dk-mock-server.yaml

