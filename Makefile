build:
	npm run build

deploy:
	 aws-vault exec icsolutions-dev \
   	npm run deploy

deploy-development:
	 aws-vault exec icsolutions-dev \
   	npm run deploy-dev

deploy-production:
	 aws-vault exec icsolutions-prod \
   	npm run deploy-prod
