<p align="center">
  <a href="https://github.com/zclzone/vue-naive-admin">
    <img alt="Vue Naive Admin Logo" width="200" src="./src/assets/images/logo.png">
  </a>
</p>
<p align="center">
  <a href="./LICENSE"><img alt="MIT License" src="https://badgen.net/github/license/zclzone/vue-naive-admin"/></a>
</p>

---

<a href="https://hellogithub.com/repository/54f19ba1f9ae4238b3cbd111f3c428b0" target="_blank"><img src="https://abroad.hellogithub.com/v1/widgets/recommend.svg?rid=54f19ba1f9ae4238b3cbd111f3c428b0&claim_uid=jXGayRdJZScqMNr" alt="Featured｜HelloGitHub" style="width: 250px; height: 54px;" width="250" height="54" /></a>

## Introduction

Vue Naive Admin is a minimalist admin template with a full-stack solution. The front-end is built with Vite + Vue 3 + Pinia + Unocss, and the back-end uses NestJS + TypeORM + MySQL. It's easy to use and visually clean. The project has gone through many refactors and polish cycles to provide a practical starter for admin interfaces.

## Design Philosophy

The Vue Naive Admin project went open-source in February 2022. From v1.0 to v2.0 it follows the principle "simplicity is virtue". It aims to help small teams, students, and solo developers quickly start admin projects. To lower the learning curve, the front-end is implemented in JavaScript (not TypeScript), making this project one of the few Vue 3 admin templates that use JavaScript while still offering a polished experience.

## Features

- 🆒 Modern Vue 3 stack: Vite + Vue 3 + Pinia
- 🍇 Atomic CSS with Unocss — elegant, lightweight, and easy to use
- 🍍 Pinia for state management with optional persistence
- 🤹 Icon system based on iconify + unocss, supports custom icons and dynamic rendering
- 🎨 Built with Naive UI for a simple and clean code style and UI; themes are easy to customize
- 👏 Clear, well-organized file structure with low coupling between modules — removing a single business module won't affect others
- 🚀 Flat route design: each component can be a page, avoiding multi-level router KeepAlive issues
- 🍒 Dynamic routes generated from permissions; 403 and 404 pages are handled separately
- 🔐 Integrates session refresh (e.g., Redis-backed) for seamless login state handling and better UX
- ✨ Global message utilities wrapped around Naive UI, supporting batch notifications and singleton cross-page behavior
- ⚡️ Common business components (Page, CRUD table, Modal, etc.) are pre-wrapped for faster development and less duplication

## Performance

![](https://docs.isme.top/Public/Uploads/2023-11-18/6558568b2b476.png)
![](https://docs.isme.top/Public/Uploads/2023-11-18/655853caa9ce8.png)

## What's new in 2.0 vs 1.0

- v2.0 is a ground-up redesign inspired by v1.0; while they may look similar, the code structure in v2.0 is significantly different.
- v1.0 provided only a front-end with mocked back-end data. v2.0 is full-stack and includes real back-end API examples.
- Although the version number is higher, v2.0 aims to be simpler and more flexible than v1.0.
- v2.0 offers higher flexibility — you can customize layouts per page if needed.

[Try v1.0 | template.isme.top](https://template.isme.top)

[Try v2.0 | admin.isme.top](https://admin.isme.top)

## NestJS Back-end

This project provides a reference back-end using NestJS + TypeORM + MySQL. It includes JWT, RBAC, and basic APIs commonly needed by the template.

- Source (GitHub): [isme-nest-serve | github](https://github.com/zclzone/isme-nest-serve)
- Source (Gitee): [isme-nest-serve | gitee](https://gitee.com/isme-admin/isme-nest-serve)

## Documentation

- Project docs: [docs | vue-naive-admin](https://isme.top)
- API docs: [apidoc | isme-nest-serve](https://apifox.com/apidoc/shared-ff4a4d32-c0d1-4caf-b0ee-6abc130f734a)

Note: A common question is how to add or modify menus. Menu resources are controlled by the back-end, so after connecting to a back-end you should use the resource management feature to add/update/delete menus, then grant permissions in the role management screen. For details about back-end integration, see the project documentation. If you prefer not to control some menus via permissions, you can add `basePermissions` in `/src/settings.js` to align with your menu resource structure (refer to the API docs for structure details).

## Get started with this template

[Create a GitHub repo from this template](https://github.com/zclzone/vue-naive-admin/generate).

Or clone using `degit` (this removes commit history):

```cmd
npx degit zclzone/vue-naive-admin
```

## License

This project is licensed under the MIT License. You are free to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the software, provided that the original copyright and license notices are preserved in the source files.

In short, the author keeps copyright but does not impose other restrictions.

## Related back-end projects

- [isme-java-serve](https://github.com/DHBin/isme-java-serve): A lightweight Java back-end using Spring Boot, MyBatisPlus, SaToken, MapStruct, etc., compatible with Vue Naive Admin 2.0.
- [naive-admin-go](https://github.com/ituserxxx/naive-admin-go): A Go back-end based on gin, gorm, MySQL, JWT, and session — compatible with Vue Naive Admin 2.0.
- [isme-java](https://github.com/AllenDengMs/isme-java): A concise Java back-end based on Spring Boot 3 + JDK21, with account and permission management, API auth, and message internationalization support.

## Contact / Community

[https://www.isme.top/contact.html](https://www.isme.top/contact.html)
