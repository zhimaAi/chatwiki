<p align="center"><a href="https://Chatwiki.com/"><img src="./imgs/head_image.png" alt="head_image"></a></p>

<p align="center">
  <a href="./README.md">English</a> |
  <a href="./README_zh.md">简体中文</a> |
  <a href="./UpdateLog.md">UpdateLog</a> |
  <a href="https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/ykeoauc4g9k2dwv1">Help Docs</a>
</p>

## 🎯 Product Positioning

ChatWiki is a workflow automation platform focused on the WeChat ecosystem, dedicated to making every official account a
super AI agent. It fully integrates the open capabilities of the official account platform, allowing you to build WeChat
ecosystem applications through drag-and-drop, enabling features like one-click rewriting of official account articles
and AI-curated comment replies.

![product positioning](./imgs/product_positioning.png)

## ✨ Core Features

### 💬 Deep Integration with WeChat Ecosystem

- **Industry First**: Automatic reply to private messages for unverified official accounts, supporting text, voice,
  images, mini-program cards, video messages, etc.

- **WeChat Workflows**: Integrates trigger scenarios such as user private messages, comments, follows, unfollows, menu
  clicks, etc. Supports various processing steps like replying to private messages, tagging fans, generating draft
  articles, publishing articles, and more.

- **Knowledge Base Synchronization**: Supports scraping articles and materials from official accounts to build a
  knowledge base with one click.

### 🤖 Basic Capabilities

- **Workflow Orchestration**: Conversational workflows, plugin workflows, including basic workflow nodes, bidirectional
  MCP, Agent mode, and user interaction.

- **Document Knowledge Base**: Supports URL reading, batch document import, API integration, AI-based segmentation, QA
  segmentation, parent-child segmentation. Supports knowledge graphs, hybrid vector search, and visual exploration of
  knowledge graphs.

- **QA Knowledge Base**: Automatically extracts QA knowledge from uploaded documents, supports automatic clustering of
  unknown questions, and summarizes common FAQs from human conversations.

- **Human Handoff**: Handles general user inquiries via bot, while also supporting human客服. Issues that the bot cannot
  resolve can be escalated to human客服, with multi-agent collaborative assignment.

- **Model Support**: Supports over 20 mainstream global models, including DeepSeek R1, doubao pro, qwen max, OpenAI,
  Claude, etc.

### 🌐 More Capabilities

- **Multiple Deployment Options**: Offers desktop client, supports publishing as a WebApp, and embedding into websites,
  official accounts/service accounts, WeChat客服, WeChat store客服, etc.

- **MCP & API Integration**: Allows integration of external MCP services or publishing workflows as MCP services.
  Complete OpenAPI interface for easy integration with existing business systems.

- **Multi-account Permission Management**: Three-tier permission system (admin, editor, viewer) ensures data isolation.
  IP whitelist and permanent login logs.

## 🛸 UI

- 🌍**Free Trial URL**: [chatwiki.com](https://chatwiki.com/)
- 🖼️**Screenshots**:

<p align="center">   <img src="./imgs/ui_1.png" alt="1" width="49%" />   <img src="./imgs/ui_2.png" alt="2" width="49%" /> </p> 
<p align="center">   <img src="./imgs/ui_3.png" alt="3" width="49%" />   <img src="./imgs/ui_4.png" alt="4" width="49%" /> </p> 

## 🚀 One-Click Deployment

ChatWiki Community Edition is deployed using Docker and can be installed in just a few simple steps:

```
# Install Docker
sudo curl -sSL https://get.docker.com/ | CHANNEL=stable sh
# Clone the project
git clone https://github.com/zhimaAi/chatwiki.git
cd chatwiki/docker
# Start the service
docker compose up -d
# Start using it, access via IP:port (ensure the specified port ${CHAT_SERVICE_PORT}, default 18080, is open)
# Default username: admin
# Default password: chatwiki.com@123
```

For any issues or suggestions during installation and deployment,
please [contact us](https://github.com/zhimaAi/chatwiki?tab=readme-ov-file#contact-us) or refer to
the [help documentation](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1?source=aHR0cHM6Ly9jaGF0d2lraS5jb20v) for
assistance. You can also check the guides below.

- [Installing ChatWiki via Installation Assistant](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/tvwn5npk63aqikq1)

- [One-Click Deploy ChatWiki Community Edition](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/wql8ekkylbwegbzo)

- [Docker Mirror Site Installation + Offline Installation](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/aa3htgexhdocyagr)

- [Deploy ChatWiki Without Docker](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/klriercbhpy97o0g)

- [Deploy ChatWiki Community Edition on Baota Linux Panel](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/gefgwdfnclua7d9y)

- [Deploy ChatWiki Community Edition with 1Panel](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/munvto5g1ctc1gcu)

- [How to Configure Model Providers and Supported Models](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/pn79lkvl53bo0xxm)

- [Local Model Deployment](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/evmy0rr9gr2gp2i0)

- [How to Configure External Services and Domain for Push Notifications](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/nfk4slc95s4i8u4v)

- [How to Obtain Large Model API Keys](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/lx3ho90skq95dpdq)

## 💻 Tech Stack

----

- Frontend: vue.js

- Backend: golang + python

- Database: PostgreSQL16 + pgvector + zhparser

<h2>🏡 Community & Contact Us <a name="contact-us"></a></h2>

----
Feel free to contact us for help or to provide suggestions to improve ChatWiki. You can reach us through:

- **Help:** Check
  the [help documentation](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1?source=aHR0cHM6Ly9jaGF0d2lraS5jb20v)
- **Email:** Send an email to [jarvis@2bai.com.cn](mailto:jarvis@2bai.com.cn) to contact us.
- **WeChat:** Scan the QR code below to join the ChatWiki tech community group. Please add the note "chatwiki" when
  adding.

<p align="left"><img src="./imgs/contact-us.png" alt="contact-us"></p>

## 📖**Changelog**

---
For the complete changelog, please click 👉️👉️[UpdateLog.md](./UpdateLog.md)

**2026/03/20**

1.Variables are displayed directly at the top of the conversation after being filled in<br/>
2.Feishu Multidimensional Table - Query Record Node: Fixed an issue where query conditions and query fields would
disappear when re-selecting the multidimensional table input box after entering data<br/>
3.Added access domain to image URLs replied by WeCom bot<br/>
4.[STD] Added ChatWiki model service to cloud model services<br/>

**2026/03/13**

1.ChatWiki supports authorized ChatClaw login<br/>
2.Knowledge Base: Support video uploads, and robots can reply with videos<br/>
3.Robot/Workflow: External services support WeCom robot<br/>
4.Fixed parameter validation error when publishing scheduled triggers<br/>
5.Workflow Publishing: Channel links displayed in the publish window<br/>
6.ChatClaw client token management and forced logout<br/>

**2026/03/06**

1. Multimodal content display compatibility in conversation records<br/>
2. [STD] International registration flow, supporting Google and email registration<br/>
3. WeChat Official Account articles - manual sync and auto sync<br/>
4. Conversation logs: Add prompt logs to conversations<br/>
5. [STD] Add universal invitation code, to be shared in community groups<br/>
6. Support merging adjacent segments when editing segments in regular knowledge base<br/>
7. Model management: Support for OpenRouter model integration<br/>
8. Core service startup optimization: Remove Neo4j dependency restriction<br/>
9. Support modifying API domain in add/edit model dialog<br/>
10. Open API: Add knowledge base recall endpoint<br/>
11. Bot/Workflow knowledge base recall metadata filtering supports referencing variables<br/>

## License

---

This project follows the [ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE).
The [ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE) is based on the Apache License
2.0, but with additional restrictions:

1. ChatWiki is free for individual users, including non-commercial or commercial activities conducted by individuals.
2. Any company, organization, institution, or team that uses ChatWiki for commercial purposes must contact us to obtain
   a commercial license.
3. When using ChatWiki's frontend components, you may not remove or modify the "ChatWiki" logo, trademark, or copyright
   notice contained therein.

**The full license text can be found in the [LICENSE](./LICENSE) file. For commercial licensing,
please [contact us](#contact-us).**
