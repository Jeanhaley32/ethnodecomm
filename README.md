# Submission for the EthGlobal 2023 Hackathon. 

## EthnodComm
 Ethnod Communicator is derived from methods created in the [go-ethereum](https://github.com/ethereum/go-ethereum) library
 Mainly those logic from [cmd/devp2p](https://github.com/ethereum/go-ethereum/tree/master/cmd/devp2p)

 - it Uses the `lookuppubkey` method with in the [`p2p/discover/v4_udp.go`](https://github.com/ethereum/go-ethereum/blob/master/p2p/discover/v4_udp.go) library To discover the neighbors of TargetNodes.

## How to use it 
 - Takes in an `enode://` encoded bootnode via the `targetnode` flag. It them derives a set of nearby neighbors
   based on that bootnodes public key.  

   Example with bootnode:
   ```
    go run main.go --targetnode enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303
   ```
   results
   ```
    enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303
enode://beaae88b39cc9baea4c3fb9af98a8a56f58b65f99cbb3fae1391819e981616f16f69ffbd864ef051fb1317a13c3065cff49feaae2fc28dc077e1d85c7589c803@79.132.239.243:30388?discport=55724
enode://4ac03ed80ea8670593c15c7bb2b46ff9d9e8bac1fe02ebd25926ddb4855dec699de044a13d6079cc160e872e7d99a0a1a962c946470df2936b5c19f7e81b1972@78.61.5.30:30303
enode://b164ae2b30e39642a3e1c22e545f92b34a7ff6047d9a3a14e99598ec17896d9f6512ce40a326fdd1574485a597d419a28d53b70de18ba324efada01c7050da6a@210.178.106.171:30303
enode://ffc50138f70f8847232a7398cf115a1edcb2334758183b5a4e6dd8057ecf41328389404f2beb0b5dbb160b02e7a75bc9a8cfb2dba0dbf970b6e4efbd1098e4e0@45.32.102.86:30405
enode://2d66797bb1cb9938df45dfd418074aae04d2802716da3fb6540d9c7173e8445b26ac80092d47d47938b2154daf31945d21e9617545f2bc1ef97218e67ea1e598@108.160.213.50:30311
enode://531c304ce89fb48bba88cdeea10d0ea5e35fdc33709683b35a16f881fb1fa50b72c3e60d89f45bdfa8e1b89c93dd6afdf2443ca58b642a24f6a803f225749ce1@129.213.99.153:30303
enode://22a1e6542679435ccc4e1551f8c35ac538a07057084d01c5416cd027fd857f91ceb2fd251d9414a5e565e852eb12029ddceb3faef752ff477a5b2c33b7d0a355@85.10.210.200:30311
enode://ea6d67eb3277d8ae9292fc700fa757ef6d2127c4db9712bcd5eb1341b1d937ac71cc2b15efe3a8496f4fc9fc12156d7ac73d82eb3c0f68928442116030b76f48@3.135.122.4:30303
enode://a5de00df98166744a64cb4a66bac3a4dac0ab86fad01d8642a0b45eacf2edf48d2a8c097b23b2e37084e8e4c143f2486cee8a6bf4c342f570450311855404103@34.91.245.204:30303
enode://7a9640265a9f248f81e533162dcd7d1a14dc4722e2b80464b269da6b279f69883985e3fd2a68411820347a9a8cc32cf0aff627b26d48c323390460e2907d73ba@217.76.53.14:39797
enode://61d057adbf2b9760900184094ef0d18985d159ad44daddcfaa0a315f45d0f8984a9883f4a23fe096da30e4c39105ea32afae3380221f0517a303aa064810e324@95.216.136.128:30303
enode://dcfc216324f8b69423b268a1d3f377780c2c51fe59b5a899ed48db9a89a74dcd2dd5f9a213eff6ab6a2f1c026dbb8df07f1097ba5b6d0c0f38701a1614feeffe@45.143.197.92:5050
enode://add782e5cb1e267203a133feaafe606225189f4cc00c8c4cf6f2cba7f2aa7a6f0a66ba62fdfb7c3fb851f5cd0ca04944bfe2a8c0c9ed603ab38e71390e4f2664@88.198.63.76:30303?discport=5742
enode://cdccced1f9664be30ca1917316be71a9ed2b746936bab1df8a66dad8b8deaf9b89f7f30b6bc83ffbde6265d1fb342c04a74c5c43f884ec8a9bec36e6044ec3c6@65.21.88.230:30303
enode://3d34d43dcfb684a2bfceb6d2f463e920e57a054d46c93db3459ce19b7ca7457e126117dfad1e576641e5341bfb554d4ee3221aa28ec26d41d9fc608f656b24bc@82.152.174.189:30311
   ```
   
## Flags
 - `targetnode` - "Target node to connect to in format `enode://` format"

