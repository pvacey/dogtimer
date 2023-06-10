all: node go

node:
	npm run build
go:
	cd backend; GOOS=linux GOARCH=arm64 go build .
tar:
	rm dogtimer | true
	rm dogtimer.tgz | true
	cp backend/dogtimer .
	tar czvf dogtimer.tgz dogtimer Caddyfile dist/
	rm dogtimer | true
distribute: node go tar
	ssh ubuntu@freeloaders.lol "rm -rf ~/dogtimer/; mkdir ~/dogtimer"
	scp dogtimer.tgz ubuntu@freeloaders.lol:~/dogtimer/
	ssh ubuntu@freeloaders.lol "cd ~/dogtimer/; tar xzvf dogtimer.tgz"
deploy: distribute
	ssh ubuntu@freeloaders.lol "cd ~/dogtimer/;pkill dogtimer;sh -c 'nohup ./dogtimer > backend.log 2>backend.error.log < /dev/null &'; sudo caddy stop; sh -c 'sudo nohup caddy run > caddy.log 2> caddy.error.log < /dev/null &'"