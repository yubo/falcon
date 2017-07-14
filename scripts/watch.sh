#!/bin/bash  
watch=./bin/*

inotifywait -mrqe close_write $watch |  while read file 
do
        module=$(basename $file)
        echo "$(date) ${module} update"
        if [ "$module" = "falcon-ctrl" ]; then
		echo "kill && start  falcon-ctrl"
                killall falcon-ctrl
                nohup ./bin/falcon-ctrl -config ./etc/falcon.conf >>log/ctrl.log 2>&1 &
        fi
done
