#!/bin/bash
cd ~/agent-mesh-core
export PATH=$PATH:$HOME/.local/go/bin:$HOME/go/bin:$HOME/.local/bin
pkill nats-server || true
pkill vextra || true
nohup nats-server -p 4222 > nats.log 2>&1 &
sleep 2
nohup ./vextra > vextra.log 2>&1 &
sleep 3
tail vextra.log
