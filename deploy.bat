@echo off
scp main.go go.mod config.json aws:~/mail
ssh aws < remote.sh