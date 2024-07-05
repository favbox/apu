# 微信公众号文章采集器

## 功能列表：
- 使用 [`ArticleCrawler.Index`](#) 采集文章列表
- 使用 [`ArticleCrawler.UpdateState`](#) 更新阅读量等扩展信息 
- 使用 [`ArticleCrawler.UpdateContent`](#) 更新内容并转为结构化数据

## 原理纪要：
获取阅读量请求权的关键是拦截到一个正常请求的 cookie（如拦截PC或手机），然后组装POST即可。