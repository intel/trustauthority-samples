# Build and Deployment of Key Broker Service and TDX Demo App

## Prerequisites
- Docker >= 20.10

## Key Broker Service
### Build
1. Download the [Dockerfile](./kbs/Dockerfile) for building KBS
   ```bash
   wget https://raw.githubusercontent.com/arvind5/trustauthority-samples/task/dockerfile-deployment/deployment/kbs/Dockerfile
   ```
2. Run the following command to build the Docker image for the Key Broker Service
   ```bash
   sudo docker build -t relying-party .
   ```

### Deploy
1. Download the [kbs.env](./kbs/kbs.env) file and update the missing configurations
   ```bash
   wget https://raw.githubusercontent.com/arvind5/trustauthority-samples/task/dockerfile-deployment/deployment/kbs/kbs.env
   ```
2. Run the following command to start the Key Broker Service as a Docker container
   ```bash
   sudo docker run --name relying-party -d --env-file kbs.env --restart=always -p 9443:9443 relying-party:latest
   ```

## TDX Demo workload
### Build
1. Download the [Dockerfile](./sample-workload/Dockerfile) for building TDX demo workload
   ```bash
   wget https://raw.githubusercontent.com/arvind5/trustauthority-samples/task/dockerfile-deployment/deployment/sample-workload/Dockerfile
   ```
2. Run the following command to build the Docker image for the TDX demo workload
   ```bash
   sudo docker build --build-arg PLATFORM=dcap -t trustauthority-demo .
   ```
   Where `PLATFORM` is the underlying TDVM where you want to run the workload, e.g., azure, gcp, dcap (Intel Dev Cloud)

### Deploy
1. Download the [workload.env](./sample-workload/workload.env) file and update the missing configurations
   ```bash
   wget https://raw.githubusercontent.com/arvind5/trustauthority-samples/task/dockerfile-deployment/deployment/sample-workload/workload.env
   ```
2. Run the following command to start the TDX demo workload as a Docker container
   ```bash
   sudo docker run --name ita-demo -d --env-file workload.env --device=/dev/tdx_guest -p 12780:12780 --user 0 trustauthority-demo:latest
   ```
   On successful run, the container will generate a key on KBS, `key_id` of which could be fetched from container logs. Make note of `key_id` from container logs to be used later in /v1/key end-point of workload.

> [!Note]
> If running on Azure, replace `--device=/dev/tdx_guest` with `--device=/dev/tpmrm0`  
> The command above will start the container as root user. This is needed because TDX demo workload container requires elevated privileges to access `/dev/tdx_guest` or `/dev/tpmrm0` device for quote collection.  
> If you do not want to run the container as root user then make sure the `/dev/tdx_guest` or `/dev/tpmrm0` device is accessible by container. One of the way to do this is to pass `gid` owning the `tdx_guest` or `tpmrm0` device to container.  
> ```bash
> sudo docker run --name ita-demo -d --env-file workload.env --device=/dev/tdx_guest -p 12780:12780 --group-add $(getent group <user-group> | cut -d: -f3) trustauthority-demo:latest
> ```
> Where `<user-group>` is the group owning the `tdx_guest` or `tpmrm0` device
> ```bash
> crw-rw---- 1 root <user-group> 10, 123 /dev/tdx_guest
> ```