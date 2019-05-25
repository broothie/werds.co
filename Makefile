
deploy: deploy-docker
	docker push gcr.io/werds-241615/werds

deploy-docker:
	docker build -t gcr.io/werds-241615/werds .
