# 假充值预防

原贴地址 <https://paper.seebug.org/853/>

针对此类型的交易，相关项目方以及交易所及钱包需要对 EOS 转账的执行状态进行校验，确保交易执行状态为 “executed”，除此之外，在区块不可逆的情况下，也要做到以下几点防止其他类型的 “假充值” 攻击的发生

- 判断 action 是否为 transfer
- 判断合约账号是否为 eosio.token 或其它 token 的官方合约账号
- 判断代币名称及精度
- 判断金额
- 判断 to 是否是自己平台的充币账号