<div align="center">
<img src="webs/src/assets/logo.png" width="150px" height="150px" />
</div>

<div align="center">
  <img src="https://img.shields.io/badge/Vue-5.0.8-brightgreen.svg"/>
  <img src="https://img.shields.io/badge/Go-1.22.2-green.svg"/>
  <img src="https://img.shields.io/badge/Element%20Plus-2.6.1-blue.svg"/>
  <img src="https://img.shields.io/badge/license-MIT-green.svg"/>
  <div align="center"> <a href="README.md">中文</a> | English</div>


</div>

# Project Introduction

`sublinkE` is an open-source node subscription conversion and generation management system based on secondary development of [sublinkX](https://github.com/gooaclok819/sublinkX).

- Frontend based on [vue3-element-admin](https://github.com/youlaitech/vue3-element-admin);
- Backend using Go + Gin + Gorm;
- Default account: admin Password: 123456, please make sure to change it after installation;

# Modifications


- [x] Fixed some BUGs
- [x] Support Clash `dialer-proxy` attribute
- [x] Allow adding and using API KEY to access API
- [x] Import and schedule updates of nodes from subscription links
- [ ] ...

# Project Features

- High flexibility and security, supporting subscription record access and simple configuration management;
- Support for multiple client protocols and formats, including:
    - v2ray (base64 universal format)
    - clash (support ss, ssr, trojan, vmess, vless, hy, hy2, tuic, AnyTLS)
    - surge (support ss, trojan, vmess, hy2, tuic)
- Added token authorization and subscription import functionality, enhancing security and convenience.

# Installation Instructions
```bash
docker run --name sublinke -p 8000:8000 \
-e API_ENCRYPTION_KEY="your-encryption-key" \  //If you need to use API functionality, please set this
-v $PWD/db:/app/db \
-v $PWD/template:/app/template \
-v $PWD/logs:/app/logs \
-d eun1e/sublinke 
```

# Project Preview

![Preview1](webs/src/assets/1.png)
![Preview2](webs/src/assets/2.png)
![Preview3](webs/src/assets/3.png)
![Preview3](webs/src/assets/4.png)