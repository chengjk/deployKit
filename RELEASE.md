
## v1.3
取消 path，url，lurl 三者之一必须传值的限制。应用场景：只想在远程服务器执行一些代码，并不上传任何东西。

## v1.2
ssh 增加 PublicKey 认证方式。
配置文件 server配置中增加 publicKey 选项。

认证顺序：

1. 有password 就用password 认证。
2. 否则检查 publicKey选项,用指定的私钥认证，
3. 如果还没有，使用用户目录 `~/.ssh/id_rsa`,
4. 如果依然没有，不能通过认证.


## v1.1
更新参数

 `-tag` 替换`-version`参数，`-scmd`替换`-cmd`;增加参数 `-pcmd`，增加参数 `-v`。

- pcmd: prefixCmd 前置命令，上传文件前执行。
- scmd：suffixCmd 后置命令，上传文件后执行。
- v : 显示版本。

更新example。

## v1.0
第一个稳定可用版。