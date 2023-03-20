build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chaxelib-local

upload:
# 需要前往服务器通过pkill -9 start停止服务，才能更新服务器
# Ipkdev2018888
	scp -r ./chaxelib-local admin@192.168.8.5:/share/SVN_Data/svn_repos/chaxelib-local/start
	ssh admin@192.168.8.5 "cd /share/SVN_Data/svn_repos/chaxelib-local"
# nohup ./start --model=product >./console.log 2>&1 &