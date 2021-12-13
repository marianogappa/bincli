## Installation

Easy version: download the [latest release](https://github.com/marianogappa/bincli/releases/latest) for your OS.

Otherwise, you know the hard version. Clone the repo, check that I'm not stealing your data and compile it.

## Usage

```
Examples:

bincli balance
bincli balance BTC
bincli chartBalanceBtc > index.html && open index.html
bincli chartBalanceUsdt > index.html && open index.html
bincli chartBalanceDataBtc
bincli chartBalanceDataUsdt
bincli ticker
bincli ticker BTCUSDT
bincli alert BTCUSDT ">" 56000 && cowsay "Reached"
bincli sushiswapAlert DELTA ">" 0.1 && cowsay "Reached"
bincli uniswapAlert APYS ">" 0.1 && cowsay "Reached"
bincli honeyswapAlert DAI ">" 0.1 && cowsay "Reached"
bincli bitforexAlert coin-usdt-omi ">" 0.1 && cowsay "Reached"
bincli kucoinAlert BTC-USDT ">" 0.1 && cowsay "Reached"
bincli bitmaxAlert BTC/USDT ">" 56000 && cowsay "Reached"
bincli ftxAlert BTC/USD ">" 56000 && cowsay "Reached"
bincli sovAlert ">" 56000 && cowsay "Reached"
bincli sushiswapTicker DELTA
bincli uniswapTicker APYS
bincli honeyswapTicker DAI
bincli bitforexTicker coin-usdt-omi
bincli kucoinTicker BTC-USDT
bincli bitmaxTicker BTC/USDT
bincli ftxTicker BTC/USD
bincli sovTicker
bincli ethGas {rapid|fast|standard|slow}
```
