SHELL := /bin/bash

# Get extra args from the command line after `make <target> .. ...`.
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

# Include .prisma.env if it exists.
-include .prisma.env
export

# Define a phony target to ensure .prisma.env exists
.PHONY: ensure-env

ensure-env:
	@if [ ! -f .prisma.env ]; then \
		echo "Copying .prisma.env.example to .prisma.env"; \
		cp .prisma.env.example .prisma.env; \
		echo "Please fill in the .prisma.env file with the correct values and rerun make."; \
		exit 1; \
	fi

clean:
	rm -rf tmp prisma/db

db: db/generate db/migrate db/spy

db/generate: ensure-env
	cd util/database/prisma && go run github.com/steebchen/prisma-client-go generate

db/migrate: ensure-env
	cd util/database/prisma && go run github.com/steebchen/prisma-client-go db push

db/spy:
	rm -rf tmp/schemaspy
	docker run --platform linux/amd64 --rm -it -v "./tmp/schemaspy":/output  --network riser schemaspy/schemaspy:latest -t pgsql -db riser -host riser-infra-timescaledb -s public -u postgres -p postgres -hq -imageformat svg
	cp tmp/schemaspy/diagrams/summary/relationships.real.large.svg docs/erd.svg
	@if [ "$(RUN_ARGS)" = "--open" ]; then open tmp/schemaspy/index.html; fi