# wxbot
wechat bot.

## 运行项目
- 配置wxbot.conf文件
- 配置文件的全路径作为参数启动应用(应用仅接受1个参数或者没有参数，不传参数则按照默认的配置启动)
- 默认配置请参考wxbot.conf

## API清单
- 支持的http请求请参考router/router.go文件
- 所有请求Header都需要带token，该token由配置文件配置
- Header:token = xxxxxx

```go
重新登陆微信接口
POST host:port/v1/wechat/login
参数类型:application/x-www-form-urlencoded
参数列表:
id = xxxxxx
id需要已经绑定了微信联系人或者群聊
返回值:提示语，用于扫码登陆的地址
```

```go
返回一个存在并没有使用过的id
GET host:port/v1/wechat/hook/id
参数:无
返回值:id
```

```go
将id绑定到需要发送消息的微信联系人或者群聊
POST host:port/v1/wechat/hook/bind
参数类型:application/x-www-form-urlencoded
参数列表
name = xxx
id = xxxxxx
name为需要绑定的微信联系人或者群聊名称，也可以是昵称，如果没有名称的话，也可以是当前登陆的id
id为/v1/wechat/hook/id接口中获取的id
返回值:成功的话将返回一个webhook地址。向该地址发送的消息都将通过机器人转发到绑定联系人或者群
```

```go
发送消息
POST host:port/v1/wechat/hook/:id/send
参数类型:application/x-www-form-urlencoded
参数列表
message = xxx
message为需要发送到微信的消息
:id为之前接口获取到的绑定的id
返回值:发送结果
```

## 使用流程
- 配置wxbot.conf文件
- 将wxbot.conf文件的全路径作为参数运行本程序，您也可以修改run.sh文件并执行
- 扫码登陆，如果远程登陆请使用浏览器打开[host:port/qr]()扫码
- 使用扫码的微信号加入群聊，并将群聊保存到通讯录中(如果不保存到通讯录中有可能会无法发送消息)
- 调用[/v1/wechat/hook/id]()接口获取发送消息id
- 调用[/v1/wechat/hook/bind]()接口将消息id与需要发送消息的联系人或者是群聊绑定并获取真实发送消息地址
- 调用上一步获取的发送消息地址发送消息，接口如[/v1/wechat/hook/:id/send]()所示


## 如有疑问请留