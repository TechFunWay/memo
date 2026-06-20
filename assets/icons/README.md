# 应用图标源文件

备忘录 App 图标的可编辑源与生成脚本。

## 文件

| 文件 | 说明 |
|------|------|
| `memo-icon.svg` | **矢量母版（编辑这个）**。emerald→teal→cyan 渐变圆角方形 + 白色便签卡片(右上折角) + 文字线 + 右下对勾徽标，含立体光影。四角透明。 |
| `master-1024.png` | 从 SVG 光栅化的 1024×1024 透明底母版（位图中转，供 sips 降采样）。 |
| `render.html` | 浏览器渲染母版用的页面（无 rsvg-convert 时用它截图）。 |
| `build-icons.sh` | 生成所有尺寸并分发到前端 / fnpack 目录。 |

## 怎么改图标

1. 编辑 `memo-icon.svg`（任何矢量/文本编辑器）。
2. 重新生成 `master-1024.png`：
   - 装了 `librsvg`：`brew install librsvg`，脚本会自动从 SVG 重生成；**或**
   - 没装：用浏览器打开 `render.html`，把 1024×1024 区域以**透明背景**截图，覆盖 `master-1024.png`。
3. 运行 `bash assets/icons/build-icons.sh`，自动生成 256/128/64/32 并分发到：
   - `frontend/public/favicon.svg` + `favicon.png`（站点 title 图标）
   - `fnpack/app/ui/images/icon_{256,128,64}.png`（fnOS 桌面图标）
   - `fnpack/ICON_256.PNG` + `ICON.PNG`（fnOS 安装包图标）
4. 要更新发布包：重跑 `./scripts/build-all.sh` 和 `./scripts/build-fnpack.sh`。

## 制作方式说明

母版是手写 SVG（矢量），用 Chromium 无头渲染（`omitBackground` 保证四角透明）光栅化成 PNG，再用 macOS `sips` 降采样到各尺寸。无第三方设计软件依赖。
