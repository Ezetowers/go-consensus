#!/bin/bash
echo 'alias ll="ls -l"' > /root/.bashrc

# Initialize serf agent
echo "Node addr: $NODE_ADDR"
serf agent -node=$NODE_ADDR -bind=$NODE_ADDR:5001 &
sleep 5

# Join the cluster
if [ ! -z "$CLUSTER_ADDR" ] ; then
    serf join $CLUSTER_ADDR:5001
fi

serf monitor