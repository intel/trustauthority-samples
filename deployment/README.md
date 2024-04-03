# Build and Deployment of Key Broker Service and TDX Sample App

## Prerequisites

- Docker >= 20.10

## ITA Key Broker Service
The Key Broker Service can be built and deployed as a Docker container. The following steps describe how to deploy the Key Broker
Service as a Docker container.

### Build
1. Download the [Dockerfile](./kbs/Dockerfile) for building KBS 
2. Run the following command to build the Docker image for the Key Broker Service  
   ```bash
   docker build -t relying-party .
   ```
    If proxy is enabled then run the following command:  
    ```bash
    docker build --build-arg HTTP_PROXY=http://<proxy_host>:<proxy_port> --build-arg HTTPS_PROXY=http://<proxy_host>:<proxy_port> -t relying-party .
    ```
### Deploy
#### Steps
1. Download the [kbs.env](./kbs/kbs.env) file and update the missing configurations
2. Run the following command to start the Key Broker Service as a Docker container with restart always
   ```bash
   docker run --name relying-party -d --env-file kbs.env  --restart=always -p 9443:9443 relying-party:latest
   ```
     
# Deploying ITA Sample workload as a container on a TD VM
### Build
1. Download the [Dockerfile](./sample-workload/Dockerfile) for building trustauthority sample workload 
2. Run the following command to build the Docker image for the ITA Sample workload image  
   ```bash
   docker build -t trustauthority-demo .
   ```
    If proxy is enabled then run the following command:  
    ```bash
    docker build --build-arg http_proxy=http://<proxy_host>:<proxy_port> --build-arg https_proxy=http://<proxy_host>:<proxy_port> --build-arg no_proxy=<KBS_IP/KBS_DNS> -t trustauthority-demo .
    ```
### Deploy
#### Steps
1. Download the [workload.env](./sample-workload/workload.env) file and update the missing configurations
2. Run either of following command to start the ITA Sample workload as a Docker container  
   i. Running as root user. The application running within container needs access to /dev/tdx_quest. Hence, it requires us to run the container as root user.  
      ```shell
      docker run --name ita-demo -d --env-file workload.env  --device=/dev/tdx_guest -p 9000:12780 --user 0  trustauthority-demo:latest
      ```
   ii. Running as non-root user.  
Make sure the Intel(R) TDX driver device is set with the following permissions:   
      ```bash
      crw-rw---- root <user-group> /dev/tdx_guest   
      ```
      Run the following command 
      ```bash
      docker run  --env-file workload.env  --device /dev/tdx_guest -p 9000:12780 --group-add $(getent group <user-group> | cut -d: -f3)  --privileged trustauthority-demo:latest   
      ```

  
