#!/bin/bash
# Copyright (C) 2024 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause

# Script is used for executing workload flow
# Runing script:
# execute_workload_flow.sh <optional_env_file>
# optional_env_file: Optional environment file path. Default is /tmp/execute_workload_flow.env

############################### Constants ###############################################################
readonly CODE_ERROR='\033[1;31m' #RED_COLOR
readonly CODE_OK='\033[0;32m'  #GREEN_COLOR
readonly TITLE='\033[1;34m' #Title/ Green
readonly CODE_NC='\033[0m' #NO_COLOR`
readonly DEFAULT_FILE_NAME="/tmp/execute_workload_flow.env"

#labels and description
readonly LABEL_GET_ATT_TOKEN="Get attestation token"
readonly LABEL_GET_ATT_TOKEN_DESC="Attesting workload with Intel Trust Authority."

readonly LABEL_GET_DESC_KEY="Get decryption key"
readonly LABEL_GET_DESC_KEY_DESC="Requesting encryption key from KBS."

readonly LABEL_DECRYPT_KEY="Decrypt model"
readonly LABEL_DECRYPT_KEY_DESC="Decrypting Model using encryption key from KBS."

readonly LABEL_EXE_MODEL="Execute model"
readonly LABEL_EXE_MODEL_DESC="Executing model with test data."

readonly LABEL_RST_MODEL="Reset Model"
readonly LABEL_RST_MODEL_DESC="Clearing decrypted model from memory."

#link
readonly HELP_LINK="https://github.com/arvind5/trustauthority-samples/blob/main/tdx-ml-workload/README.md"
##############################################################################################
function checkIfRequireToolInstalled() {
    loc="$(type -p "$1")"
    if ! [ -f "$loc" ]; then
        return 1
    fi
    return 0
}

function checkIfFileExists() {
    if [ ! -f "$1" ]; then
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

function validateEnvironmentFile(){
    envFileName=$DEFAULT_FILE_NAME
    #Read input variable for environment file. If not exist use default file
    if [ $# -eq 1 ];then
        envFileName=$1    
    fi
    printf "Using $envFileName environment file\n"
    
    #check if environment file exists
    checkIfFileExists $envFileName
    if [ $? == 1 ]; then
        echo -e "\u274c  Environment file ($envFileName) does not exist, exiting... Please create env file and rerun script"
        return 1
    fi
    export $(grep -v '^#' $envFileName | xargs -d '\n')
    
    hasError=false
    if [[ -z "$WORKLOAD_URL" ]];then
        printf "${CODE_ERROR}%s${CODE_NC}" "WORKLOAD_URL must be present"
        hasError=true
    fi    
    if [[ -z "$KBS_URL" ]];then
        printf "${CODE_ERROR}%s${CODE_NC}" "KBS_URL must be present"
        hasError=true
    fi  
    if [[ -z $KBS_KEY_ID ]];then
        printf "${CODE_ERROR}%s${CODE_NC}" "KBS_KEY_ID must be present"
        hasError=true
    fi  
    if [[ $hasError == true ]];then
        printf "${CODE_ERROR} %s ${CODE_NC}" "Exiting..."
        return 1
    fi  
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

printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =
#check and validate environment file
validateEnvironmentFile $1
if [ $? == 1 ]; then
    exit 1
fi

#labels
label_url="Url"
label_method="Method"
label_response="Response"
label_request="Request"
label_expected_code="Expected Code"
label_received_code="Received Code"
label_decoded_token="Decoded token"

server_addr=$WORKLOAD_URL
kbs_host=$KBS_URL
key_id=$KBS_KEY_ID

printf "Using Workload Host url : $server_addr\n"
printf "Using KBS host url      : $kbs_host\n" 
printf "Using KBS Key id        : $key_id\n"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' =

#Variable for returning response
return_var=''
#Variable for returning code
return_code=''
#Get token
url=$server_addr/taa/v1/token

printf "\n\n"
printf "${TITLE} ${LABEL_GET_ATT_TOKEN}: ${CODE_NC}${LABEL_GET_ATT_TOKEN_DESC}\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| GET"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 200"
makeHttpsCall $url GET return_var return_code
if [ $return_code -ne 200 ]; then
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf "%*s %s\n" ${#label_expected_code} "$label_response" "| $return_var"

attestationToken=$(jq -r -n --argjson data "$return_var" '$data.attestation_token')
printf "\n\n${CODE_OK}%*s %s${CODE_NC}\n" ${#label_expected_code} "$label_decoded_token"
printf "\n\n${CODE_OK}%*s %s${CODE_NC}\n" ${#label_expected_code} " $(jq -R 'split(".") | .[1] | @base64d | fromjson' <<< \"$attestationToken\")"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -

#Make key call
requestForKey=$( jq -n \
                  --arg at "$attestationToken" \
                  --arg kturl "$kbs_host/keys/$key_id/transfer" \
                  '{attestation_token: $at, key_transfer_url: $kturl}' )
url=$server_addr/taa/v1/key

printf "\n${TITLE} ${LABEL_GET_DESC_KEY}: ${CODE_NC}${LABEL_GET_DESC_KEY_DESC}\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 200"
printf "%*s %s\n" ${#label_expected_code} "$label_request" "| $requestForKey"

makeHttpsCall $url POST return_var return_code "$requestForKey"
if [ $return_code -ne 200 ]; then
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf "%*s %s\n" ${#label_expected_code} "$label_response" "| $return_var"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -

#make decrypt call
url=$server_addr/taa/v1/decrypt
printf "\n${TITLE} ${LABEL_DECRYPT_KEY}: ${CODE_NC}${LABEL_DECRYPT_KEY_DESC}\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 204"
printf "%*s %s\n" ${#label_expected_code} "$label_request" "| $return_var"

makeHttpsCall $url POST return_var return_code "$return_var"
if [ $return_code -ne 204 ]; then
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_response" "| $return_var"
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
printf "\n${TITLE} ${LABEL_EXE_MODEL}: ${CODE_NC}${LABEL_EXE_MODEL_DESC}\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 200"
printf "%*s %s\n" ${#label_expected_code} "$label_request" "| $requestForExecute"

makeHttpsCall $url POST return_var return_code "$requestForExecute"
if [ $return_code -ne 200 ]; then
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf "%*s %s\n" ${#label_expected_code} "$label_response" "| $return_var"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -

#reset model
url=$server_addr/taa/v1/reset
printf "\n${TITLE} ${LABEL_RST_MODEL}: ${CODE_NC}${LABEL_RST_MODEL_DESC}\n\n"
sleep 5
printf "%*s %s\n" ${#label_expected_code} "$label_url" "| $url"
printf "%*s %s\n" ${#label_expected_code} "$label_method" "| POST"
printf "%*s %s\n" ${#label_expected_code} "$label_expected_code" "| 204"

makeHttpsCall $url POST return_var return_code
if [ $return_code -ne 204 ]; then
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_received_code" "| Received unexpected status $return_code, exiting..."
    printf "%*s ${CODE_ERROR}%s${CODE_NC}\n" ${#label_expected_code} "$label_response" "| $return_var"
    exit 1
fi
printf "%*s %s\n" ${#label_expected_code} "$label_received_code" "| $return_code"
printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -