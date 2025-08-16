#!/bin/bash

# Presumed environment variables, ensure these are created in provisioner
# SCRIPT_DBPASS
# SCRIPT_DBUSER
# SCRIPT_DBNAME
# SCRIPT_DBROOTPASS
# SCRIPT_DUMPPATH
# SCRIPT_WEBSERVERDOCKERPATH

DBPORT=5433

echo "Installing Docker..."
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Upgrade system packages
sudo apt upgrade -y

#Setup Docker Network
sudo docker network create --subnet=172.18.0.0/16 net

echo "Installing Postgresql..."
# Get PostgreSQL docker image
sudo apt-get install -y postgresql-client 
sudo docker pull postgres
# Setup data
mkdir data
CONTAINERID=$(sudo docker run -d --restart=always --net net --ip 172.18.0.2 -e POSTGRES_PASSWORD=$SCRIPT_DBROOTPASS -v $(realpath data):/var/lib/postgresql/data -p $DBPORT:5432 postgres:latest)
sleep 10 # Give time for the postgres instance to finish booting.
PGPASSWORD=$SCRIPT_DBROOTPASS psql --host=localhost -p $DBPORT -U postgres -f $SCRIPT_DUMPPATH
rm $SCRIPT_DUMPPATH
sudo apt-get remove -y postgresql-client

#Install Portfolio Webserver image
echo "Installing PortfolioWebserver..."
sudo docker load -i $SCRIPT_WEBSERVERDOCKERPATH

sudo docker run -d --restart=always --net net --ip 172.18.0.3 -p 80:80 -e PORTFOLIOSERVER_DBIP=172.18.0.2 -e PORTFOLIOSERVER_DBUSER=$SCRIPT_DBUSER -e PORTFOLIOSERVER_DBPORT=5432 -e PORTFOLIOSERVER_DBPASS=$SCRIPT_DBPASS -e PORTFOLIOSERVER_DBNAME=$SCRIPT_DBNAME -e PORTFOLIOSERVER_PORT=80 portfoliowebserver:latest 

echo "Setting up firewall..."
#Enable firewall, close off SSH
# sudo ufw allow 80
# sudo ufw --force enable

#Connection should be terminated, ready to be saved as an AMI