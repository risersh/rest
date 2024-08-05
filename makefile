# Include .prisma.env if it exists.
-include .prisma.env

ensure-env:
ifeq ("$(wildcard .prisma.env)","")
	echo "Copying .prisma.env.example to .prisma.env"; \
	cp .prisma.env.example .prisma.env; \
	echo "Please fill in the .prisma.env file with the correct values and rerun make."; \
	exit 1;
endif

clean:
	rm -rf tmp prisma/db

db/generate:
	@if [ ! -f .prisma.env ]; then \
		echo "Copying .prisma.env.example to .prisma.env"; \
		cp .prisma.env.example .prisma.env; \
		echo "Please fill in the .prisma.env file with the correct values and rerun make."; \
		exit 1; \
	fi
	cd prisma && go run github.com/steebchen/prisma-client-go generate

db/migrate: ensure-env
	cd prisma && go run github.com/steebchen/prisma-client-go db push

db/spy:
	rm -rf tmp/schemaspy
	docker run --platform linux/amd64 --rm -it -v "./tmp/schemaspy":/output  --network riser schemaspy/schemaspy:latest -t pgsql -db riser -host riser-infra-timescaledb -s public -u postgres -p postgres -hq -imageformat svg
	if command -v open; then open tmp/schemaspy/index.html; fi
