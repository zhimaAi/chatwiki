<p align="center"><a href="https://Chatwiki.com/"><img src="./imgs/head_image.png" alt="head_image"></a></p>

<p align="center">
  <a href="./README_en.md">English</a> |
  <a href="./README.md">简体中文</a> |
  <a href="./UpdateLog.md">更新日志</a> |
  <a href="https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/ykeoauc4g9k2dwv1">帮助文档</a>
</p>

## 🎯 产品定位

ChatWiki 是一个专注微信生态的工作流自动化平台，致力于让每个公众号都可成为一个超级AI智能体。全面集成公众号平台的开放能力，拖拽即可搭建微信生态应用，实现公众号推文一键改稿、留言AI精选回复等能力

![product positioning](./imgs/product_positioning.png)

## ✨ 核心特性

### 💬 微信生态深度集成

- **全行业首创**：未认证公众号私信自动回复，支持文本、语音、图片、小程序卡片、视频消息等。

- **微信工作流**：集成用户私信、留言、关注、取关、点击菜单等触发场景，支持回复私信，粉丝打标签，生成草稿文章、发布文章等多种处理流程

- **知识库同步**：支持抓取公众号文章素材，一键建立知识库。

### 🤖 基础能力

- **工作流编排：** 对话工作流、插件工作流，包含基础的工作流节点、双向 MCP、Agent 模式、用户交互。

- **文档知识库：** 支持 url 读取、文档批量导入、API 对接、支持AI分段、QA分段、父子分段。支持知识图谱、向量混合检索，可视化查看知识图谱。

- **问答知识库：** 上传文档自动抽取问答知识，支持未知问题自动聚类，支持从人工对话中总结常用FAQ

- **转人工客服：** 通过机器人处理一般的用户咨询，同时支持人工客服接待。机器人处理不好的问题可以由人工客服介入处理，支持多客服协同分配。

- **模型支持：** 支持DeepSeek R1、doubao pro、qwen max、Openai、Claude 等全球20多种主流模型。

### 🌐 更多能力

- **多种部署方式**：提供桌面客户端、支持发布为WebApp，支持嵌入网站、公众号服务号、微信客服、微信小店客服等

- **MCP&API集成**：可引入外部MCP服务，或将工作流发布为MCP服务。完整的OpenAPI接口，轻松集成现有业务系统。

- **多账号权限管理**：管理、编辑、查看三级权限体系，实现数据权限隔离。IP白名单、登录日志永久留存。

## 🛸UI

- 🌍**免费体验网址**： [chatwiki.com](https://chatwiki.com/)
- 🖼️**系统截图**：

<p align="center">   <img src="./imgs/ui_1.png" alt="1" width="49%" />   <img src="./imgs/ui_2.png" alt="2" width="49%" /> </p> 
<p align="center">   <img src="./imgs/ui_3.png" alt="3" width="49%" />   <img src="./imgs/ui_4.png" alt="4" width="49%" /> </p> 

## 🚀 一键部署

ChatWiki 社区版基于 Docker 部署，只需简单几步即可完成安装：

```
# 安装 Docker
sudo curl -sSL https://get.docker.com/ | CHANNEL=stable sh
# 克隆项目
git clone https://github.com/zhimaAi/chatwiki.git
cd chatwiki/docker
# 启动服务
docker compose up -d
# 开始使用，通过IP+端口访问(需要开放指定的端口${CHAT_SERVICE_PORT},默认18080)
# 默认账号：admin
# 默认密码：chatwiki.com@123
```

在安装和部署中有任何问题或者建议，可以 [联系我们](https://github.com/zhimaAi/chatwiki?tab=readme-ov-file#contact-us)
或者查看 [帮助文档](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1?source=aHR0cHM6Ly9jaGF0d2lraS5jb20v)
获取帮助，也可以参考下面的文档。

- [通过chatwiki安装助手安装](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/tvwn5npk63aqikq1)

- [一键部署ChatWiki社区版](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/wql8ekkylbwegbzo)

- [docker镜像站安装+离线安装](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/aa3htgexhdocyagr)

- [免Docker部署ChatWiki](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/klriercbhpy97o0g)

- [宝塔Linux面板部署ChatWiki社区版](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/gefgwdfnclua7d9y)

- [使用1Panel部署ChatWiki社区版](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/munvto5g1ctc1gcu)

- [如何配置模型供应商及支持的模型](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/pn79lkvl53bo0xxm)

- [本地模型部署](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/evmy0rr9gr2gp2i0)

- [如何配置对外服务和接收推送的域名](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/nfk4slc95s4i8u4v)

- [如何获取大模型ApiKey](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/lx3ho90skq95dpdq)

## 💻 技术栈

----

- 前端：vue.js

- 后端：golang +python

- 数据库：PostgreSQL16+pgvector+zhparser

<h2>🏡社区交流&联系我们 <a name="contact-us"></a></h2>

----
欢迎联系我们获取帮助，或者提供建议帮助我们改善ChatWiki。您可以通过以下方式联系我们：

- **帮助：** 查看 [帮助文档](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1?source=aHR0cHM6Ly9jaGF0d2lraS5jb20v)
- **邮箱：** 您可以发送邮件到 [jarvis@2bai.com.cn](mailto:jarvis@2bai.com.cn)联系我们。
- **微信：** 使用微信扫码加入ChatWiki技术交流群，添加请备注“chatwiki”

<p align="left"><img src="./imgs/contact-us.png" alt="contact-us"></p>

## 📖**更新日志**

---
查看完整更新日志请点击👉️👉️[UpdateLog.md](./UpdateLog.md)

**2026/01/16**

1.后台字体默认显示为平方字体<br/>
2.新增博查联网搜索插件<br/>
3.增加webhook触发器以及结束节点支持输出变量<br/>
4.工作流:工作流支持调用工作流<br/>
5.【STD】添加会话转人工筛选和导出功能<br/>
6.【STD】机器人:会话记录显示转人工的标识支持筛选<br/>
7.飞书多维表查询记录条件等支持输入Json格式的内容<br/>
8.扩展TokenUseReport回调函数以支持图像数量参数<br/>
9.机器人菜单问题修复,恢复显示功能中心,未知问题<br/>
10.新增插件:探索>插件广场>百度搜索插件<br/>
11.【STD】模型管理:云版chatwiki自带模型增加图片生成模型<br/>
12.修复循环节点的输出key赋值为中间变量key的问题<br/>
13.聊天机器人支持清空聊天缓存<br/>
14.飞书多维表增加添加协作者<br/>
15.工作流:http节点支持增加鉴权参数<br/>
16.前端页面去掉版权信息<br/>
17.机器人应用卡片支持拖动排序<br/>
18.工作流:新增立即回复节点,支持立即输出消息<br/>
19.普通知识库:支持导入飞书知识库<br/>
20.【STD】邀请注册流程调整,增加用户交流群<br/>
21.工作流-数据库缺失自动修复<br/>
22.新增插件:探索>插件广场>OCR工具解析<br/>
23.仅显示头像不显示账号信息<br/>
24.兼容麒麟操作系统python多线程运行没权限<br/>
25.工作流:http节点支持单独测试，支持自动提取输出参数<br/>
26.【STD】兼容麒麟系统离线激活机器码生成失败处理<br/>

**2026/01/09**

1.工作流:节点卡片增加渐变色<br/>
2.工作流:模板/机器人文字内容超长,鼠标移入时显示全部<br/>
3.修复获取指定机器人功能报错问题<br/>
4.工作流:节点前端效果优化,触发器备注字段不美观问题修复<br/>
5.工作流:新增应用收费功能<br/>
6.修复开始节点变量无法删除问题<br/>
7.前端小问题合集+前端优化合集<br/>
8.应用收费海报图片跨域修复<br/>
9.公众号文章群发插件、留言管理-插件<br/>
10.工作流:增加获取公众号accesstoken的节点<br/>
11.公众号粉丝数据、图文群发数据<br/>
12.修复模型被删除 选择模型组件无法选择问题<br/>
13.工作流:JSON序列化和反序列化节点类型支持<br/>
14.模板广场:模板支持上传使用说明<br/>
15.工作流:增加知识库导入节点<br/>
16.微信公众号文章爬虫集成<br/>
17.rerank重排输入取全量的rrf的结果<br/>
18.工作流插件状态异常修复<br/>
19.工作流:执行超时调整到60分钟,测试时显示时长与token消耗<br/>
20.【STD】注册页面添加邀请码字段并优化样式<br/>
21.线上功能体验问题整理、工作流-测试时参数给默认值<br/>
22.多语言国际化配置(系统管理模块和登陆页面)<br/>
23.公众号触发器变量显示异常问题修复<br/>
24.【STD】关键词回复和收到消息回复-支持小客服消息API接入机器人场景<br/>
25.修复公众号触发器保存后未返回消息字段问题<br/>
26.父子分段:同一个父分段下的子分段用虚线分隔<br/>
27.多语言设置框调整到账号下拉中<br/>

**2025/12/31**

1.触发器输出支持不配置变量映射+补充开始节点自定义全局变量<br/>
2.【STD】开放接口:/v1/chat/completions 兼容多模态输入<br/>
3.工作流:注释卡片节点改成弹窗打开填写<br/>
4.工作流:新增语音合成节点和声音复刻节点<br/>
5.公众号草稿:获取草稿详情和插件安装后即开启/飞书插件的兼容<br/>
6.新增插件:探索>插件广场>公众号智能接口<br/>
7.爬虫服务(crawler):抓取隐藏识别标识<br/>
8.线上工作流问题和优化合集<br/>
9.工作流结束节点支持自定义回复消息<br/>
10.H5App对外服务中，工作流返回较慢时，前端会不停重发消息<br/>
11.【STD】商业版增加离线激活<br/>
12.添加插件初始化枚举变量和修复动态模版引用<br/>
13.聊天请求流式输出增加定时keep-alive消息保活逻辑<br/>
14.飞书多维表支持创建多维表/创建数据表/创建视图<br/>
15.rerank流程改为串行以及RRF算法增加加权数<br/>
16.工作流:http节点支持导入CURL<br/>

## 协议

---

本项目遵循[ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE)
开源协议。[ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE)基于Apache License
2.0协议，但是有一些额外的限制：

1. ChatWiki 对个人用户免费，包括个人从事的非商业或商业活动。
2. 任何公司、组织、机构或团队若将 ChatWiki 用于商业目的，均须联系我们获得商业授权。
3. 在使用 ChatWiki 的前端组件时，您不得移除或修改其中包含的“ChatWiki”标识、商标或版权声明。

**完整的许可证文本请查看：[LICENSE](./LICENSE) 文件，需要获取商业授权请[联系我们](#contact-us)**

