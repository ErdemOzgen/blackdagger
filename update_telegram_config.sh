#!/bin/bash

# Usage check
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <TELEGRAM_API_KEY> <TELEGRAM_CHAT_ID>"
    exit 1
fi

# Set the file path
file_path="/root/.config/notify/provider-config.yaml"
new_api_key=$1
new_chat_id=$2

# Assuming the YAML structure is very predictable and simple
# This is a fragile solution and might not work for all YAML structures
sed -i "/id: \"tel\"/,/telegram_api_key:/s/telegram_api_key:.*/telegram_api_key: ${new_api_key}/" "$file_path"
sed -i "/id: \"tel\"/,/telegram_chat_id:/s/telegram_chat_id:.*/telegram_chat_id: ${new_chat_id}/" "$file_path"
