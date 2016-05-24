# ccms
基于beego框架封装的通用api接口框架<br>
Features： <br>
1、接口版本化<br>
2、接口安全签名验证<br>
3、入参统一处理<br>
4、参数xss过滤<br>
5、利用vendor管理第三方包的版本统一性<br>
6、mysql、redis封装及配置化<br>

#项目run起来后，本地测试链接
http://localhost:8080/v2/person/get/?appVersion=3.1.0&platform=ios&requestId=592385239&signature=157541d3b4a43675f00a3cf3707603697dea73d2&timestamp=1458289047&userId=10001
