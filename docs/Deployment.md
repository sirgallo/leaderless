# Deployment

This implementation includes `docker-compose` configuration for running a cluster locally. The `compose` file will run a cluster of 5 athn nodes, as well as a forward facing `haproxy` instance as the loadbalancer to send requests to.

First, ensure `docker engine` and `docker compose` are installed on your system (for `macos`, this involves installing `docker desktop`). [Click Here](https://www.docker.com/products/docker-desktop/) to download the latest version of `docker desktop`.

The basic implementation to run the cluster and the associated `docker` resources are located under [cmd](./cmd)

Once `docker desktop` is installed, run the following to deploy locally (on `macos`):

  1. Make sure HOSTNAME is set in ~/.zshrc and registered in /etc/hosts

```bash
export HOSTNAME=hostname
source ~/.zshrc
```

once sourced, restart your terminal for the changes to take effect

```bash
sudo nano /etc/hosts
```

open the file and add:
```bash
127.0.0.1 <your-hostname>
```

this binds your hostname to localhost now so you can use your hostname as the address to send commands to

  2. Generate the self signed certs

Haproxy is served over https, so generate your self signed certs under the [certs](./certs/) folder in the root of the project. This can be done with

```bash
cd ./certs
openssl req -newkey rsa:2048 -new -x509 -days 365 -nodes -out $HOSTNAME.crt -keyout $HOSTNAME.key
cat $HOSTNAME.key $HOSTNAME.crt > $HOSTNAME.pem
```

`docker-compose` will then bind the certs from your local machine to the certs folder on the haproxy container. Your hostname will also be bound as the hostname of the container

  3. Run the startupDev.sh script

Again, this is located in the root of the project. Run the following:

```bash
chmod +x ./startupDev.sh
./startupDev.sh
```

This will build the services and then start them. 

At this point, you can begin interacting with the cluster by sending it commands to perform on the state machine. 

To stop the cluster, run the `./stopDev.sh` script:
```bash
chmod +x ./stopDev.sh
./stopDev.sh
```