# 整洁架构 Clean Architecture

## 关键点
- 框架隔离
- 数据独立
- 业务集中
- 随时部署

## 参考文献
- [The Clean Code Blog by Robert C. Martin (Uncle Bob)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Proposal：工程化模板或标准化的讨论](https://github.com/cloudwego/kitex/issues/500)
- [Clean Architecture – Caching As A Proxy](https://codecoach.co.nz/clean-architecture-caching-as-a-proxy/)
- [Clean Architecture, 2 years later](https://eltonminetto.dev/en/post/2020-07-06-clean-architecture-2years-later/)


# go包的组织

## 包的类别 package category: lib + app
- app = business domain model + service(use case -> data)

- [2023 Clean Architecture](https://medium.com/inside-picpay/organizing-projects-and-defining-names-in-go-7f0eab45375d)
- [2020 Clean Architecture](https://eltonminetto.dev/en/post/2020-07-06-clean-architecture-2years-later/)


## 库包 packages/library
- 包：只用于单一目的
- 子包：用于相同目的的不同方法
- 包名：要求简短精确

## 应用包 packages/application
- 快速测试
- 易于理解
- 易于重构
- 易于维护

