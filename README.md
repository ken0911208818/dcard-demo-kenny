# 練習Dcard的面試題 
### 參考 https://github.com/KennyChenFight/dcard-simple-demo
### 啟動 
```shell script
docker-compose up -d
```

### redis 筆記
> 筆記如何redis 實現 transaction 的過程 並表示 redis 原生不支援 transatction 必須要透過Lua 的腳本才能完成

>> https://hackmd.io/vWUOZrxTQzC8q10r-LHCvw

### gin validate 
> 待補

### orm 
作者的orm是使用`xorm` 本專案中我使用的是 `gorm` 

### validate v10 

在v10 的版本中 許多 func 都需要在最前面加入`ctx` `c := ctx.background()`

