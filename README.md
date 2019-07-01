# go-oneday


#### 项目介绍
golang自定义web框架，作为学习掌握Golang编写的实践。
1.支持restful以及占位符参数
2.支持流式注册接口
3.支持设置全局http，以及单独设置某个接口的http

#### 软件架构
>odserver 框架入口
>>httpConfig http配置类
>>Router 路由控制器

>router
>>request封装客户端请求
>>pathMatcher 路径匹配工具
>>HandlerObject 封装请求信息
>>FuncObject 请求对应的接口函数

>request 自定义请求实体类


#### 安装教程

#### 使用说明
```
func main() {
	o := odserver.Default()
	o.Start("/main").
		Target("/test/").Get(HelloServer).Post(HelloServer).Delete(HelloServer).
		And().
		Target("/test2").Get(HelloServer2)

	o.Start("/{test}/main/").Target("/number/{number}").
		Get(HelloServer3).Post(HelloServer4)

	http.ListenAndServe(":8080",o)
}
```

#### 参与贡献


