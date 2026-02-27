<p align="center"><a href="https://Chatwiki.com/"><img src="./imgs/head_image.png" alt="head_image"></a></p>

<p align="center">
  <a href="./README_en.md">English</a> |
  <a href="./README.md">简体中文</a> |
  <a href="./UpdateLog.md">UpdateLog</a> |
  <a href="https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/ykeoauc4g9k2dwv1">Help Docs</a>
</p>

## 🎯 Product Positioning

ChatWiki is a workflow automation platform focused on the WeChat ecosystem, dedicated to making every official account a super AI agent. It fully integrates the open capabilities of the official account platform, allowing you to build WeChat ecosystem applications through drag-and-drop, enabling features like one-click rewriting of official account articles and AI-curated comment replies.

![product positioning](./imgs/product_positioning.png)

## ✨ Core Features

### 💬 Deep Integration with WeChat Ecosystem

- **Industry First**: Automatic reply to private messages for unverified official accounts, supporting text, voice, images, mini-program cards, video messages, etc.

- **WeChat Workflows**: Integrates trigger scenarios such as user private messages, comments, follows, unfollows, menu clicks, etc. Supports various processing steps like replying to private messages, tagging fans, generating draft articles, publishing articles, and more.

- **Knowledge Base Synchronization**: Supports scraping articles and materials from official accounts to build a knowledge base with one click.

### 🤖 Basic Capabilities

- **Workflow Orchestration**: Conversational workflows, plugin workflows, including basic workflow nodes, bidirectional MCP, Agent mode, and user interaction.

- **Document Knowledge Base**: Supports URL reading, batch document import, API integration, AI-based segmentation, QA segmentation, parent-child segmentation. Supports knowledge graphs, hybrid vector search, and visual exploration of knowledge graphs.

- **QA Knowledge Base**: Automatically extracts QA knowledge from uploaded documents, supports automatic clustering of unknown questions, and summarizes common FAQs from human conversations.

- **Human Handoff**: Handles general user inquiries via bot, while also supporting human客服. Issues that the bot cannot resolve can be escalated to human客服, with multi-agent collaborative assignment.

- **Model Support**: Supports over 20 mainstream global models, including DeepSeek R1, doubao pro, qwen max, OpenAI, Claude, etc.

### 🌐 More Capabilities

- **Multiple Deployment Options**: Offers desktop client, supports publishing as a WebApp, and embedding into websites, official accounts/service accounts, WeChat客服, WeChat store客服, etc.

- **MCP & API Integration**: Allows integration of external MCP services or publishing workflows as MCP services. Complete OpenAPI interface for easy integration with existing business systems.

- **Multi-account Permission Management**: Three-tier permission system (admin, editor, viewer) ensures data isolation. IP whitelist and permanent login logs.

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

For any issues or suggestions during installation and deployment, please [contact us](https://github.com/zhimaAi/chatwiki?tab=readme-ov-file#contact-us) or refer to the [help documentation](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1?source=aHR0cHM6Ly9jaGF0d2lraS5jb20v) for assistance. You can also check the guides below.

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

- **Help:** Check the [help documentation](https://www.yuque.com/zhimaxiaoshiwangluo/pggco1?source=aHR0cHM6Ly9jaGF0d2lraS5jb20v)
- **Email:** Send an email to [jarvis@2bai.com.cn](mailto:jarvis@2bai.com.cn) to contact us.
- **WeChat:** Scan the QR code below to join the ChatWiki tech community group. Please add the note "chatwiki" when adding.

<p align="left"><img src="./imgs/contact-us.png" alt="contact-us"></p>

## 📖**Changelog**

---
For the complete changelog, please click 👉️👉️[UpdateLog.md](./UpdateLog.md)

**2026/02/06**

1. Optimized thumbnail generation for knowledge base files.<br/>
2. [STD] Real-time conversation: Supports batch handoff to human客服.<br/>
3. Fixed issue where types from Excel imports were not correctly handled when deleting QA knowledge bases.<br/>
4. [STD] Custom models: Support for Tongyi Qianwen/Wanxiang image models.<br/>
5. Knowledge base: Word document parsing now supports formula import.<br/>
6. [STD] Official website: Added official account AI agent tab.<br/>
7. New plugin: DingTalk multi-dimensional table.<br/>
8. Added privacy policy consent reminder.<br/>
9. Workflow: Automatically locates the error card when testing/publishing validation fails.<br/>
10. Trigger count: Statistics module filters out deleted content.<br/>
11. Workflow: Code node now supports running Python code.<br/>

**2026/01/30**

1. Workflow: Variable notes are displayed in process start/run test/logs.<br/>
2. Documentation/Search menus are closed by default.<br/>
3. [STD] Added recycle bin for QA knowledge bases.<br/>
4. Feishu multi-dimensional table - Add data/Update record nodes support passing variables.<br/>
5. [STD] Knowledge base recycle bin supports batch delete/restore.<br/>
6. Knowledge base: Added unknown question statistics, supports viewing details.<br/>
7. External services: WebApp and website support starting new conversations.<br/>
8. When exporting CSL and uploading templates, plugin authentication fields are cleared.<br/>
9. Workflow: Added HTTP tool node.<br/>
10. Fixed error caused by duplicate unique key when initializing triggers.<br/>
11. [STD] Fixed sorting field error in robot list.<br/>
12. Workflow LLM capability node supports setting prompt character roles.<br/>
13. Workflow: Added QA node, supports responding to user queries.<br/>
14. Knowledge base supports switching to folder view.<br/>
15. [STD] Custom models support Gemini, GPT, Claude.<br/>
16. Application: Chatbot/Workflow supports retrieving adjacent segments.<br/>
17. Frontend page optimization: View invitation code, supports one-click copy.<br/>
18. Fixed issue where bot replies still occurred after disassociating an official account from a bot.<br/>
19. Session list needs to support filtering by customer's entry time for consultation.<br/>
20. Image models - Support Tongyi Qianwen/Wanxiang.<br/>
21. Automatically add version number when image loading fails.<br/>
22. [STD] Commercial version: Bot quantity adjusted from 20 to 100.<br/>

**2026/01/23**

1. Knowledge base: Trigger count statistics now include knowledge base grouping levels.<br/>
2. Fixed to_tsquery syntax error in full-text search ts_rank usage.<br/>
3. Added login detection when pulling official account articles and separate login for cloud version users.<br/>
4. Adjusted parsing strategy for HTML files in knowledge base documents.<br/>
5. [STD] Cloud version: Added insufficient points reminder.<br/>
6. [STD] Cloud version: Fixed issue where sub-members need invitation code to log in.<br/>
7. Template square: Templates now include main image.<br/>
8. When exporting workflows, added check for whether knowledge base and database should be exported.<br/>
9. Workflow: Feishu node and some node inputs now have enlarged editing box.<br/>
10. Knowledge base: Added API for deleting knowledge base documents/Q&A.<br/>
11. Frontend interface: Phase 2 of multilingual translation.<br/>
12. Workflow: Run test and run logs now include input and output nodes.<br/>
13. OpenAPI: Fixed issue with gin's ShouldBind parsing structs with any type.<br/>
14. [STD] Workflow: When adding an application, display large template images.<br/>
15. Workflow: Authentication fields defined in code are automatically cleared when exporting CSL files and uploading templates.<br/>
16. External services: Verified WeChat official accounts/mini programs support replying with "Content generated by AI" during "typing..." status.<br/>
17. Knowledge base: Added metadata support and metadata filtering.<br/>
18. Disable plugin functionality and prompt when unable to access cloud version backend.<br/>
19. Robot: Added "generating answer" prompt, supports customization.<br/>
20. Robot: Supports adding variables.<br/>
21. Fixed null pointer issue in plugin status check.<br/>

## License

---

This project follows the [ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE). The [ChatWiki Open Source License](https://github.com/zhimaAi/chatwiki/blob/main/LICENSE) is based on the Apache License 2.0, but with additional restrictions:

1. ChatWiki is free for individual users, including non-commercial or commercial activities conducted by individuals.
2. Any company, organization, institution, or team that uses ChatWiki for commercial purposes must contact us to obtain a commercial license.
3. When using ChatWiki's frontend components, you may not remove or modify the "ChatWiki" logo, trademark, or copyright notice contained therein.

**The full license text can be found in the [LICENSE](./LICENSE) file. For commercial licensing, please [contact us](#contact-us).**
