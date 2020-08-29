build:
	rm isucon9final
	go build

deploy:
	scp isucon9final isucon9-a: