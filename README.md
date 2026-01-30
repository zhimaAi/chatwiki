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

**2026/01/30**

1.工作流:流程开始/运行测试/调用日志中显示变量备注<br/>
2.文档/搜索菜单默认关闭<br/>
3.【STD】问答知识库增加回收站<br/>
4.飞书多维表-新增数据/更新记录节点支持传入变量<br/>
5.【STD】知识库回收站支持批量操作删除/还原<br/>
6.知识库增加未知问题统计支持查看详情<br/>
7.对外服务:webapp、网站支持新增对话<br/>
8.导出csl和上架模板时,清空插件鉴权字段<br/>
9.工作流:新增http工具节点<br/>
10.修复触发器初始化时重复插入配置唯一键报错问题<br/>
11.【STD】修复机器人列表排序字段错误问题<br/>
12.工作流大模型能力节点支持设置提示词所属角色<br/>
13.工作流:新增问答节点,支持对用户提问<br/>
14.知识库支持切换为文件夹视图<br/>
15.【STD】自有模型支持 Gemini、gpt、Claude<br/>
16.应用:聊天机器人/工作流支持召回相邻分段<br/>
17.前端页面优化,查看邀请码,支持一键复制<br/>
18.修复机器人取消关联公众号后,还是走机器人回复<br/>
19.会话列表需要支持按照客户进入咨询的时间筛选<br/>
20.图像模型-支持通义千问/万像<br/>
21.图片加载失败自动添加版本号<br/>
22.【STD】商业版,机器人数量调整20->100<br/>

**2026/01/23**

1.知识库:触发次数统计新增知识库分组层级<br/>
2.知识库全文检索ts_rank时使用的to_tsquery语法错误修复<br/>
3.拉取公众号文章增加登录检测以及云版用户单独登录<br/>
4.知识库文档html文件解析策略调整<br/>
5.【STD】云版增加积分不足提醒<br/>
6.【STD】云版修复子成员登录需要邀请码<br/>
7.模板广场:模板增加主图<br/>
8.导出工作流,增加对知识库和数据库是否导出的判断<br/>
9.工作流:飞书节点、部分节点输入内容增加放大编辑框<br/>
10.知识库:增加知识库文档/问答删除api<br/>
11.前端界面:多语言翻译二期<br/>
12.工作流:运行测试与运行日志新增输入输出节点<br/>
13.开放接口:修复gin的ShouldBind解析结构体定义any类型报错问题<br/>
14.【STD】工作流:新增应用时,显示模板大图<br/>
15.工作流:代码中定义鉴权字段,导出csl文件上架模板时自动清空<br/>
16.对外服务:已认证微信公众号/小程序,支持回复“显示内容由AI生成“正在输入中”<br/>
17.知识库:知识库增加元数据、支持元数据过滤<br/>
18.无法访问云版后台时关闭插件功能并提示<br/>
19.机器人:增加答案生成中提示语,支持自定义<br/>
20.机器人:支持添加变量<br/>
21.修复插件状态检查空指针问题<br/>

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

## 协议

---

本项目遵循[ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE)
开源协议。[ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE)基于Apache License
2.0协议，但是有一些额外的限制：

1. ChatWiki 对个人用户免费，包括个人从事的非商业或商业活动。
2. 任何公司、组织、机构或团队若将 ChatWiki 用于商业目的，均须联系我们获得商业授权。
3. 在使用 ChatWiki 的前端组件时，您不得移除或修改其中包含的“ChatWiki”标识、商标或版权声明。

**完整的许可证文本请查看：[LICENSE](./LICENSE) 文件，需要获取商业授权请[联系我们](#contact-us)**

