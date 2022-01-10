# CHANGELOG

## v1.0.28

* 新增飞书、钉钉、企业微信机器人通知接口。

## v1.0.27

* 修复 JSON 参数类型报错。

## v1.0.26

* HttpClient 新增 Cookie Jar 支持。
* HttpClient 传入 XML, JSON 参数支持 string, []byte 类型。
* 新增 alphaID 数字、字符串转换编码。

## v1.0.23

* 新增 SearchFile 函数。

## v1.0.22

* 新增 Base58, Base64 编/解码。
* 新增 BKDR/SDBM/RS/JS/PJW/ELF/DJB/AP 哈希函数。
* HttpRequest 新增 Text 字段，便于文本协议数据 POST/PUT 需求场景。（例如: Influx 文本协议）

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
