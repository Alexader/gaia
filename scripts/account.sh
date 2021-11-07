set -e

CURRENT_PATH=$(pwd)
OPT=$1

function printHelp() {
  print_blue "Usage:  "
  echo "  account <OPT>"
  echo "    <OPT> - one of 'list', 'show'"
  echo "      - 'list' - list all account address in local testnet"
  echo "      - 'show' - show one account address by account name"
  echo "  account.sh -h (print this message)"
}

function list() {
  for ((i = 1; i < 11; i = i + 1)); do
      gaiad keys show acc${i} -a --keyring-backend test
  done
}

function show() {
  gaiad keys show acc$2 -a --keyring-backend test
}

if [ "$OPT" == "list" ]; then
  list
elif [ "$OPT" == "show" ]; then
  show
else
  printHelp
  exit 1
fi
