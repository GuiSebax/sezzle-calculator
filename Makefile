.PHONY: up down test coverage

up:
	docker compose up --build

down:
	docker compose down

test:
	$(MAKE) -C backend test
	cd frontend && npm test

coverage:
	$(MAKE) -C backend coverage
	cd frontend && npm run coverage
