build:
	go build

deploy-daemon:
	cat isutrain.service | ssh isucon@isucon9-a "sudo tee /etc/systemd/system/isutrain.service"
	ssh isucon9-a sudo systemctl daemon-reload
	cat nginx.conf | ssh isucon@isucon9-a "sudo tee /etc/nginx/nginx.conf"
	ssh isucon9-a sudo systemctl restart nginx

deploy:
	ssh isucon9-a sudo systemctl stop isutrain
	scp isucon9final isucon9-a:/home/isucon/isucon9-final/webapp/go/isucon9final
	scp start.sh isucon9-a:/home/isucon/start.sh
	ssh isucon9-a sudo systemctl start isutrain

bench:
	ssh isucon9-bench "cd isucon9-final && bench/bin/bench_linux run --payment=http://10.146.15.196:15000 --target=http://10.146.15.196:80 --assetdir=webapp/frontend/dist"

pprof:
	go tool pprof -http="127.0.0.1:8020" logs/latest/cpu.pprof

synclogs:
	rsync -av isucon@isucon9-a:/tmp/isucon/logs/ logs/
