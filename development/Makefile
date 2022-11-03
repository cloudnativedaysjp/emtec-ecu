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
	sed -e "s/__DATE__/$(shell date "+%Y-%m-%d/")" \
		-e "s/__TIME0__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "1 minute")+09:00"/" \
		-e "s/__TIME1__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "2 minute")+09:00"/" \
		-e "s/__TIME2__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "4 minute")+09:00"/" \
		-e "s/__TIME3__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "6 minute")+09:00"/" \
		-e "s/__TIME4__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "8 minute")+09:00"/" \
		-e "s/__TIME5__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "10 minute")+09:00"/" \
		-e "s/__TIME6__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "12 minute")+09:00"/" \
		-e "s/__TIME7__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "14 minute")+09:00"/" \
		-e "s/__TIME8__/"2000-01-01T$(shell date "+%H:%M:%S.%3N" -d "16 minute")+09:00"/" \
		dk-mock-server-base.yaml > dk-mock-server.yaml

