# 微信公众号文章采集器

## 功能列表：
- 使用 [`weixin.GetArticles`](#) 采集文章列表
- 使用 [`weixin.GetArticleStat`](#) 采集阅读量 
- 使用 [`weixin.GetArticleByURL`](#) 获取文章内容

## 原理纪要：
获取阅读量请求权的关键是拦截到一个正常请求的 cookie（如拦截PC或手机），然后组装POST即可。