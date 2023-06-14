remote_target = ubuntu@freeloaders.lol

all: node go

node:
	npm run build
go: node
	rm -rf backend/dist/ | true
	cp -R dist/ backend/
	cd backend; go build .
	mv backend/dogtimer backend/dogtimer_x86_64
	cd backend; GOOS=linux GOARCH=arm64 go build .
	mv backend/dogtimer backend/dogtimer_arm64
tar:
	rm dogtimer_arm64 dogtimer.tgz | true
	cp backend/dogtimer_arm64 .
	tar czvf dogtimer.tgz dogtimer_arm64 Caddyfile
distribute: node go tar
	ssh $(remote_target) "cp ~/dogtimer/times.db ~/"
	ssh $(remote_target) "rm -rf ~/dogtimer/; mkdir ~/dogtimer"
	ssh $(remote_target) "cp ~/times.db ~/dogtimer/"
	scp dogtimer.tgz $(remote_target):~/dogtimer/
	ssh $(remote_target) "cd ~/dogtimer/; tar xzvf dogtimer.tgz"
deploy: distribute
	ssh $(remote_target) "cd ~/dogtimer/;pkill dogtimer;sh -c 'nohup ./dogtimer_arm64 > backend.log 2>backend.error.log < /dev/null &'; sudo caddy stop; sh -c 'sudo nohup caddy run > caddy.log 2> caddy.error.log < /dev/null &'"