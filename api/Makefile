build:
	GOOS=linux go build -o app
	docker build -t tts_server .
	rm -f app

push:
	gcloud docker -- push gcr.io/justforfunc-prep/say
