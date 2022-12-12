#!/bin/bash
sudo docker kill --signal=SIGTERM saturn-node
echo 'wait for 30 minutes to drain all requests'
sleep 1800
sudo docker stop saturn-node
echo 'stop docker saturn ok'

