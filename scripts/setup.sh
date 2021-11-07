# You can run all of these commands from your home directory
cd $HOME
N=$1
# Initialize the genesis.json file that will help you to bootstrap the network
gaiad init my-node --chain-id my-chain

# Create a key to hold your validator account
for ((i = 1; i < N+1; i = i + 1)); do
    gaiad keys add acc${i} --keyring-backend test
    # Add that key into the genesis.app_state.accounts array in the genesis file
    # NOTE: this command lets you set the number of coins. Make sure this account has some coins
    # with the genesis.app_state.staking.params.bond_denom denom, the default is staking
    gaiad add-genesis-account $(gaiad keys show acc${i} -a --keyring-backend test) 1000000000stake,1000000000validatortoken
    # Generate the transaction that creates your validator
    gaiad gentx acc${i} 1000000000stake --chain-id my-chain --keyring-backend test
done

# generate some target address for transfer
for ((i = N+1; i < N+N+1; i = i + 1)); do
    gaiad keys add acc${i} --keyring-backend test
done

# Add the generated bonding transaction to the genesis file
gaiad collect-gentxs

# Now its safe to start `gaiad`
gaiad start