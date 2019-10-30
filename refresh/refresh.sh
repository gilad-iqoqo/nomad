set -x 
set -e

pkill -9 nomad || true
rm -rf /Users/gilad/tmp/nomad/data
SL=4

./bin/nomad agent -dev --config ~/iqoqo/nomad-example/standalone-server-conf.hcl &
sleep $SL
for i in `seq 0 3` ; do 
	echo $i
	./bin/nomad agent -dev --config ~/iqoqo/nomad-example/standalone-agent-conf-${i}.hcl | tee /tmp/client${i} &
	sleep $SL
done





