#!/bin/sh

# skip if genesis file already exists
if [ -f /root/.gaiad/config/genesis.json ]; then
  exec /entrypoint.sh
  exit 0
fi

# initialize chain
/gaiad init --chain-id localgaia local

# create keys
cat <<EOF | /gaiad keys --keyring-backend file add master --recover
dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog fossil
password
password
EOF

# create genesis accounts
/gaiad genesis add-genesis-account cosmos1hnyy4gp5tgarpg3xu6w5cw4zsyphx2lyq6u60y 10000000uatom        # validator
/gaiad genesis add-genesis-account cosmos1cyyzpxplxdzkeea7kwsydadg87357qnalx9dqz 150000000000000uatom # smoke contrib
/gaiad genesis add-genesis-account cosmos1phaxpevm5wecex2jyaqty2a4v02qj7qmhq3xz0 150000000000000uatom # smoke master
/gaiad genesis add-genesis-account cosmos1f4l5dlqhaujgkxxqmug4stfvmvt58vx2fqfdej 1000000000000uatom   # master

# replace stake with uatom
sed -i 's/"stake"/"uatom"/g' /root/.gaia/config/genesis.json

# create genesis transaction
echo "password" | /gaiad genesis gentx --keyring-backend=file master 10000000uatom --chain-id=localgaia
/gaiad genesis collect-gentxs

# enable api
sed -i '/\[api\]/,/^###/ s/^enable = false/enable = true/' /root/.gaia/config/app.toml

exec /entrypoint.sh --minimum-gas-prices=0.001uatom --grpc.address=0.0.0.0:9090 --api.address=tcp://0.0.0.0:1317
