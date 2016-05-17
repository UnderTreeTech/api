# ccms
基于beego框架封装的通用api接口框架
Features：
1、接口版本化
2、接口安全签名验证
3、入参统一处理
4、利用vendor管理第三方包的版本统一性
5、mysql、redis封装及配置化

#项目run起来后，本地测试链接
http://localhost:8080/v2/person/get/?appVersion=3.1.0&platform=ios&requestId=592385239&signature=157541d3b4a43675f00a3cf3707603697dea73d2&timestamp=1458289047&userId=10001
