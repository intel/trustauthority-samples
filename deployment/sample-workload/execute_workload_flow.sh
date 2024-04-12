#!/bin/bash

function checkIfRequireToolInstalled() {
    loc="$(type -p "$1")"
    if ! [ -f "$loc" ]; then
        return 1
    fi
    return 0
}

function makeHttpsCall(){
    #api address
    addr=$1

    #required field from json response
    method=$2

    #body
    body=$5

    response=''
    
    if [ -z "$body" ] 
    then
        response=$(curl -X $method -ksH "Accept: application/json" -w "%{http_code}" $addr)
    else
        response=$(curl -X $method -ksH "Accept: application/json" -H "Content-Type: application/json" -w "%{http_code}" -d "$body" $addr)        
    fi
    
    http_code=$(tail -n1 <<< "$response")  # get the last line
    content=$(sed '$ d' <<< "$response")  # get all but the last line which contains the status code
    export $3="$content"
    eval "$4=$http_code"  
}


#Prerequisite check for curl, jq
checkIfRequireToolInstalled "curl"
if [ $? == 1 ]; then
    echo -e "\u274c  curl is not installed, exiting... Please install curl and rerun script" 
    exit 1   
else 
    echo -e "\u2714   curl is installed"    
fi

checkIfRequireToolInstalled "jq"
if [ $? == 1 ]; then
    echo -e "\u274c  jq is not installed, exiting... Please install jq and rerun script"    
    exit 1
else 
    echo -e "\u2714   jq is installed"    
fi

#labels
label_url="Url"
label_method="Method"
label_response="Response"
label_request="Request"
label_expected_code="Expected Code"
label_received_code="Received Code"

#Required server address
echo -n "Enter Workload Host url (https://<server-ip>:<port>) :"
read server_addr
echo -n "Enter KBS host url (https://<server-ip>:<port>) :" 
read kbs_host
echo -n "Enter KBS Key id (Check trustauthority-demo container logs for key_id):"
read key_id

#Variable for returning response
return_var=''
#Variable for returning code
return_code=''
#Get token
url=$server_addr/taa/v1/token

printf "\n\n"
printf "\033[1;34m Get attestation token: \033[00mAttesting workload with Intel Trust Authority.\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| GET"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 200"
makeHttpsCall $url GET return_var return_code
if [ $return_code -ne 200 ]; then
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf "%*s %s\n" ${#label_expected_code} "$label_response" "| $return_var"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -

#Make key call
attestationToken=$(jq -r -n --argjson data "$return_var" '$data.attestation_token')
requestForKey=$( jq -n \
                  --arg at "$attestationToken" \
                  --arg kturl "$kbs_host/kbs/v1/keys/$key_id/transfer" \
                  '{attestation_token: $at, key_transfer_url: $kturl}' )
url=$server_addr/taa/v1/key
printf "\n\033[1;34m Get decryption key: \033[00mRequesting encryption key from KBS.\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 200"
printf "%*s %s\n" ${#label_expected_code} "$label_request" "| $requestForKey"

makeHttpsCall $url POST return_var return_code "$requestForKey"
if [ $return_code -ne 200 ]; then
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf "%*s %s\n" ${#label_expected_code} "$label_response" "| $return_var"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -

#make decrypt call
url=$server_addr/taa/v1/decrypt
printf "\n\033[1;34m Decrypt model: \033[00mDecrypting Model using encryption key from KBS.\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 204"
printf "%*s %s\n" ${#label_expected_code} "$label_request" "| $return_var"

makeHttpsCall $url POST return_var return_code "$return_var"
if [ $return_code -ne 204 ]; then
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -

#make execute call
requestForExecute='{
    "pregnancies": "1",
    "blood-glucose": "128",
    "blood-pressure": "88",
    "skin-thickness": "39",
    "insulin": "110",
    "bmi": "36.5",
    "dbf": "1.057",
    "age": "37"
}'
url=$server_addr/taa/v1/execute
printf "\n\033[1;34m Execute model: \033[00mExecuting model with test data.\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 200"
printf "%*s %s\n" ${#label_expected_code} "$label_request" "| $requestForExecute"

makeHttpsCall $url POST return_var return_code "$requestForExecute"
if [ $return_code -ne 200 ]; then
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf "%*s %s\n" ${#label_expected_code} "$label_response" "| $return_var"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -

#reset model
url=$server_addr/taa/v1/reset
printf "\n\033[1;34m Reset Model: \033[00mClearing decrypted model from memory.\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 204"

makeHttpsCall $url POST return_var return_code
if [ $return_code -ne 204 ]; then
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s \033[1;31m%s\033[00m\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -