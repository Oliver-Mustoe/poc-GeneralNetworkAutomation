#!/bin/bash

# Cleanup all files on server for quick testing, I scp all of the files over so I cleanup my files (rsync in future)
cd ~/poc-GeneralNetworkAutomation
sudo docker compose down
rm -rf /opt/forgejo
cd ~
rm -rf ./poc-GeneralNetworkAutomation
