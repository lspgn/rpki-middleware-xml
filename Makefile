.PHONY: clean
clean:
	rm -rf ssl
	rm -rf data1 data2

.PHONY: build-certs
build-certs:
	mkdir ssl
	openssl req -newkey rsa:2048 -nodes -keyout ssl/key1.pem -x509 -days 365 -out ssl/cert1.pem -subj "/CN=krill1"
	openssl req -newkey rsa:2048 -nodes -keyout ssl/key2.pem -x509 -days 365 -out ssl/cert2.pem -subj "/CN=krill2"
	openssl req -newkey rsa:2048 -nodes -keyout ssl/keym.pem -x509 -days 365 -out ssl/certm.pem -subj "/CN=middleware"
	cat ssl/cert*.pem >> ssl/ca-certificates.crt
