all:
	cd cmd/consumer && \
	GOOS=linux GOARCH=amd64 go build  && \
	docker build -t otel-demo-consumer-go . && \
	cd ../provider && \
	GOOS=linux GOARCH=amd64 go build && \
	docker build -t otel-demo-provider-go .