#!/bin/bash

# Cleanup all files on server for quick testing, I scp all of the files over so I cleanup my files (rsync in future)
cd ~
rm -rf poc-GeneralNetworkAutomation
rm -rf /opt/forgejo
