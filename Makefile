build:
	npm run build

synth:
	 aws-vault exec icsolutions-dev \
   	npm run synth

deploy:
	 aws-vault exec icsolutions-dev \
   	npm run deploy

deploy-development:
	 aws-vault exec icsolutions-dev \
   	npm run deploy-dev

deploy-production:
	 aws-vault exec icsolutions-prod \
   	npm run deploy-prod

run: synth
	 aws-vault exec icsolutions-dev -- sam local start-api -t ./cdk.out/soldev-api.template.json

migration-up:
	migrate -database ${POSTGRESQL_URL} -path src/migrations up

migration-down:
	migrate -database ${POSTGRESQL_URL} -path src/migrations down
