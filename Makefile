build:
	go build -o mackerel-plugin-oci-mds ./cmd/mds/...
	go build -o mackerel-plugin-oci-flb ./cmd/flb/...
	go build -o mackerel-plugin-oci-nlb ./cmd/nlb/...
