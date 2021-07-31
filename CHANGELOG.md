# CHANGELOG

## v1.0.22

* 新增 Base58, Base64 编解码。

## v1.0.21

* SubProcess 新增 GetCommand() 方法。
* SubProcess PrintCommands() 方法改为 logrus.Info 输出。
* GetDatabaseSchemas() 新增过滤 performance_schema 系统表。
* 新增 GetUsers(), GetUserPrivileges() 方法，方便查询 MySQL 用户及权限信息。

## v1.0.20

* SSH Tunnel 本地端口支持传 `0` 值。(即随机端口号)
* SubProcess 新增 ClearCommand() 方法。

## v1.0.19

* 新增 SSH Tunnel 连接支持。

## v1.0.5

* 调整 HTTPOption* 命名。

## v1.0.4

* SSHClient 新增 HTTP/Socks5 代理支持。

## v1.0.3

* Add functions: VerifySum, RemoveAllSafe
* Code optimization.
