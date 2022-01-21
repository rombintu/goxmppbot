build:
	go build main.go
run:
	go run main.go
prod: 
	go build main.go
	scp {main,config.toml} nick@mvdbot:/home/nick/opt/mvdbot
prodbuild:
	go build main.go
	cp config.toml.bak mvdbot_amd64/etc/mvdbot/config.toml
	mv main mvdbot_amd64/usr/local/bin/mvdbot
	dpkg-deb --build --root-owner-group mvdbot_amd64
	# scp mvdbot_amd64.dev nick@mvdbot:/home/nick/