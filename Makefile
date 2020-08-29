build:
	rm isucon9final
	go build

deploy:

	cat isutrain.service | ssh isucon@isucon9-a "sudo tee /etc/systemd/system/isutrain.service"
	ssh isucon9-a sudo systemctl daemon-reload

	cat nginx.conf | ssh isucon@isucon9-a "sudo tee /etc/nginx/nginx.conf"
	ssh isucon9-a sudo systemctl restart nginx

	ssh isucon9-a sudo systemctl stop isutrain
	scp isucon9final isucon9-a:/home/isucon/isucon9-final/webapp/go/isucon9final
	ssh isucon9-a sudo systemctl start isutrain

