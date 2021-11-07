# You can run all of these commands from your home directory
cd $HOME

# Initialize the genesis.json file that will help you to bootstrap the network
gaiad init my-node --chain-id my-chain

# Create a key to hold your validator account
gaiad keys add acc1 --keyring-backend test
gaiad keys add acc2 --keyring-backend test

# Add that key into the genesis.app_state.accounts array in the genesis file
# NOTE: this command lets you set the number of coins. Make sure this account has some coins
# with the genesis.app_state.staking.params.bond_denom denom, the default is staking
gaiad add-genesis-account $(gaiad keys show acc1 -a --keyring-backend test) 1000000000stake,1000000000validatortoken

# Generate the transaction that creates your validator
gaiad gentx acc1 1000000000stake --chain-id my-chain --keyring-backend test

# Add the generated bonding transaction to the genesis file
gaiad collect-gentxs

# Now its safe to start `gaiad`
gaiad start

gaiad query bank balances cosmos1uxcv2ux5jyp9grrsvgte4udfkrqfs8j26w0sg3 --chain-id my-chain