# wxbot
wechat bot.

## 运行项目
main.go 文件中包含启动一个微信机器人客户端和打开本地8888端口的http server，您可以修改该文件

## API清单
支持的http请求请参考router/router.go文件
```go
重新登陆微信接口
POST host:8888/v1/wechat/login
参数类型:application/x-www-form-urlencoded
参数列表:
token = xxxxxx
token需要已经绑定了微信联系人或者群聊
返回值:提示语，用于扫码登陆的地址

```
```go
生成一个token并将该token返回
POST host:8888/v1/wechat/hook/key
Header:token = xxxxxx
该token为代码中写死的，本来是用于页面通过API获取绑定token使用，目的是仅让管理员生成token
参数:无
返回值:token

--------------------------------------------------------------------------------
返回一个存在并没有使用过的token，如果池中没有则生成一个并返回
GET host:8888/v1/wechat/hook/key
Header:token = xxxxxx
该token为代码中写死的，本来是用于页面通过API获取绑定token使用，目的是仅让管理员生成token
参数:无
返回值:token
```

```go
将token绑定到需要发送消息的微信联系人或者群聊
POST host:8888/v1/wechat/hook/bind
参数类型:application/x-www-form-urlencoded
参数列表
name = xxx
token = xxxxxx
name为需要绑定的微信联系人或者群聊名称，也可以是昵称，如果没有名称的话，也可以是当前登陆的id
token为/v1/wechat/hook/key接口中获取的token
返回值:成功的话将返回一个webhook地址。向该地址发送的消息都将通过机器人转发到绑定联系人或者群
```

```go
发送消息
POST host:8888/v1/wechat/hook/:token/send
参数类型:application/x-www-form-urlencoded
参数列表
message = xxx
message为需要发送到微信的消息
:token为之前接口获取到的绑定的token
返回值:无
```

## 使用流程
- 运行本程序
- 扫码登陆，如果远程登陆请使用浏览器打开[host:8889/qr]()扫码
- 使用扫码的微信号加入群聊，并将群聊保存到通讯录中(如果不保存到通讯录中有可能会无法发送消息)
- 调用[/v1/wechat/hook/key]()接口获取token
- 调用[/v1/wechat/hook/bind]()接口将token与需要发送消息的联系人或者是群聊绑定并获取真实发送消息地址
- 调用上一步获取的发送消息地址发送消息，接口如[/v1/wechat/hook/:token/send]()所示


## 如有疑问请留