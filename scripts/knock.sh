#!/bin/bash
# Get path to script in your system.
SCRIPT_PATH="`dirname \"$0\"`"                    # relative
SCRIPT_PATH="`( cd \"${SCRIPT_PATH}\" && pwd )`"  # absolute
if [[ -z "$SCRIPT_PATH" ]] ; then
  exit 1  # fail, check permissions
fi

# Get path to repo in your system.
REPO_PATH=${SCRIPT_PATH%/*}

# Get path to where you are executing this script from.
EXEC_PATH=$(pwd)

cd ${REPO_PATH}

# Build the binary, if it doesn't exist in this repo.
BINARY_PATH="$REPO_PATH/port-knock"
if [[ ! -f ${BINARY_PATH} ]] ; then
  echo "binary not found for \`port-knock\`. building binary at ${BINARY_PATH}..."
  go build
fi

$(./port-knock)

cd ${EXEC_PATH}

# Uncomment me and fill in the variables with your info.
#ssh -p <ssh-port> -i <private-ssh-key> <user>@<host>
