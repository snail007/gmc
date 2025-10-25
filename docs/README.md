# GMC Framework Documentation

Welcome to the GMC Framework! Here you'll find complete documentation in both English and Chinese.

## 📖 Documentation Versions

### English Documentation
- **Status**: Core features and basic usage
- **Content**: 1,621 lines covering essential functionality
- **Location**: [MANUAL.md](MANUAL.md)
- **Online**: https://snail007.github.io/gmc/

### Chinese Documentation (中文文档)
- **Status**: Complete and detailed documentation
- **Content**: 7,479 lines with comprehensive guides
- **Location**: [zh/MANUAL_ZH.md](zh/MANUAL_ZH.md)
- **Online**: https://snail007.github.io/gmc/zh/

> **Note**: For complete and detailed documentation, we recommend reading the Chinese version. The English version provides core features and basic usage.

## 🚀 Quick Links

- [English Documentation](https://snail007.github.io/gmc/)
- [中文文档](https://snail007.github.io/gmc/zh/)
- [GitHub Repository](https://github.com/snail007/gmc)
- [API Reference](https://pkg.go.dev/github.com/snail007/gmc)

## 📚 Local Preview

### Method 1: Using docsify-cli (Recommended)

```bash
# Install docsify-cli
npm install -g docsify-cli

# Run in docs directory
cd docs
docsify serve .

# Visit http://localhost:3000
```

### Method 2: Using serve.sh

```bash
cd docs
./serve.sh

# Visit http://localhost:3000
```

### Method 3: Using Python

```bash
cd docs
python3 -m http.server 3000

# Visit http://localhost:3000
```

## 📂 Documentation Structure

```
docs/
├── index.html          # English documentation entry
├── MANUAL.md           # English complete manual (1,621 lines)
├── _sidebar.md         # English sidebar navigation
├── zh/                 # Chinese documentation
│   ├── index.html      # Chinese documentation entry
│   ├── MANUAL_ZH.md    # Chinese complete manual (7,479 lines)
│   └── _sidebar.md     # Chinese sidebar navigation
├── serve.sh            # Local preview script
└── README.md           # This file
```

## 📋 Documentation Content

### English Documentation (MANUAL.md)

Core features and basic usage guide:

1. GMC Introduction
2. Quick Start
3. Controller
4. HTTP Router
5. Template
6. Web Server
7. API Server
8. Database
9. Cache
10. I18n
11. Remote Debugging
12. Configuration File
13. Middleware
14. Official Middleware
15. Hot Update & Restart
16. Useful GMC Packages
17. Misc
18. GMCT Tool Chain

### Chinese Documentation (zh/MANUAL_ZH.md)

Complete documentation with 22 chapters:

1. GMC 框架介绍 (GMC Introduction)
2. 快速开始 (Quick Start)
3. 核心概念 (Core Concepts)
4. 配置 (Configuration)
5. 路由 (Routing)
6. 控制器 (Controller)
7. 请求与响应 (Request & Response)
8. 视图模板 (View Templates)
9. 数据库 (Database)
10. 缓存 (Cache)
11. Session
12. 日志 (Logging)
13. 国际化 (I18n)
14. API 开发 (API Development)
15. 测试 (Testing)
16. 部署 (Deployment)
17. GMCT 工具链 (GMCT Tool Chain)
18. 进阶主题 (Advanced Topics)
19. 最佳实践 (Best Practices)
20. 常见问题 (FAQ)
21. 常用工具包 (Useful Utilities)
22. 总结 (Summary)

## 🌐 Language Switching

- In the documentation, click **🌐 中文版本** or **🌐 English Version** in the sidebar
- Or visit directly:
  - English: https://snail007.github.io/gmc/
  - Chinese: https://snail007.github.io/gmc/zh/

## 🤝 Contributing to Documentation

We welcome contributions to improve the documentation!

1. Fork this project
2. Create your branch (`git checkout -b improve-docs`)
3. Commit your changes (`git commit -am 'Improve docs: xxx'`)
4. Push to the branch (`git push origin improve-docs`)
5. Create a Pull Request

## 📝 Documentation Guidelines

- Use Markdown format
- Provide complete, runnable code examples
- Keep Chinese and English documentation synchronized when possible
- Use clear heading hierarchy
- Add necessary code comments
- Test all code examples before submitting

## 💬 Feedback and Suggestions

If you find errors in the documentation or have suggestions for improvement:

- Submit an [Issue](https://github.com/snail007/gmc/issues)
- Create a [Pull Request](https://github.com/snail007/gmc/pulls)
- Join our community discussions

## 📄 License

Documentation is licensed under the [MIT License](../LICENSE)

---

**Thank you for using the GMC Framework!**

感谢您使用 GMC 框架！
