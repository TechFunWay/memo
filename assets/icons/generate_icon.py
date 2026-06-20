#!/usr/bin/env python3
"""生成备忘录 3D（等距投影）应用图标 SVG。
修改本脚本可调整图标几何/光影，运行后输出 memo-icon.svg。"""
import os

# 顶面菱形（等距盒子的顶面）顶点：N 后、E 右、S 前、W 左
# 放大占满图标，加厚以增强立体感
N = (128, 48)
E = (224, 104)
S = (128, 160)
W = (32, 104)
H = 48  # 厚度（页堆高度）

def drop(p, h=H):
    return (p[0], p[1] + h)

# 顶面局部坐标：point(a,b) = N + a*(E-N) + b*(W-N)，a,b∈[0,1] 覆盖整个菱形
EX = (E[0] - N[0], E[1] - N[1])
WX = (W[0] - N[0], W[1] - N[1])
def P(a, b):
    return (N[0] + a * EX[0] + b * WX[0], N[1] + a * EX[1] + b * WX[1])

def poly(pts):
    return " ".join(f"{x:.1f},{y:.1f}" for x, y in pts)

def quad_line(a0, a1, b0, db):
    """顶面上沿 u 方向(下右)的一条文字线（平行四边形）"""
    return [P(a0, b0), P(a1, b0), P(a1, b0 + db), P(a0, b0 + db)]

Ep, Sp, Wp = drop(E), drop(S), drop(W)
top   = [N, E, S, W]
right = [S, E, Ep, Sp]   # 前右侧面（较暗）
left  = [W, S, Sp, Wp]   # 前左侧面（中等）

# 文字线（顶面）
title = quad_line(0.20, 0.52, 0.20, 0.085)
body = [quad_line(0.22, 0.82, 0.40, 0.065),
        quad_line(0.22, 0.82, 0.56, 0.065),
        quad_line(0.22, 0.62, 0.72, 0.065)]
chk = [P(0.11, 0.27), P(0.17, 0.33), P(0.28, 0.16)]  # 勾选折线

# 侧面顶部“最上一页”的亮边
right_edge = [S, E, (E[0], E[1] + 5), (S[0], S[1] + 5)]
left_edge  = [W, S, (S[0], S[1] + 5), (W[0], W[1] + 5)]

# 侧面页堆分隔线（实现“多张纸叠起来”的立体感）
def side_line(a, b, t):
    return (a[0], a[1] + t * H), (b[0], b[1] + t * H)
stack = []
for t in (0.42, 0.72):
    (x1, y1), (x2, y2) = side_line(S, E, t)   # 右面
    stack.append(f'<line x1="{x1:.1f}" y1="{y1:.1f}" x2="{x2:.1f}" y2="{y2:.1f}" stroke="#9aa8b1" stroke-width="2.2" stroke-opacity="0.8"/>')
    (x1, y1), (x2, y2) = side_line(W, S, t)   # 左面
    stack.append(f'<line x1="{x1:.1f}" y1="{y1:.1f}" x2="{x2:.1f}" y2="{y2:.1f}" stroke="#b3bec5" stroke-width="2.2" stroke-opacity="0.8"/>')
stack_svg = "\n  ".join(stack)

# 球形对勾徽标（放大，前右悬浮）
BX, BY, BR = 204, 150, 30

svg = f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" width="256" height="256">
  <!-- 备忘录 3D 等距图标 · 由 generate_icon.py 生成，勿手改坐标 -->
  <defs>
    <linearGradient id="bg" x1="18%" y1="0%" x2="86%" y2="100%">
      <stop offset="0%" stop-color="#3BE0A0"/><stop offset="50%" stop-color="#14B8A6"/><stop offset="100%" stop-color="#0C84AE"/>
    </linearGradient>
    <radialGradient id="topGloss" cx="50%" cy="-8%" r="72%">
      <stop offset="0%" stop-color="#fff" stop-opacity="0.40"/><stop offset="55%" stop-color="#fff" stop-opacity="0.06"/><stop offset="100%" stop-color="#fff" stop-opacity="0"/>
    </radialGradient>
    <radialGradient id="botShade" cx="50%" cy="112%" r="70%">
      <stop offset="0%" stop-color="#053b3a" stop-opacity="0.30"/><stop offset="55%" stop-color="#053b3a" stop-opacity="0.05"/><stop offset="100%" stop-color="#053b3a" stop-opacity="0"/>
    </radialGradient>
    <linearGradient id="faceTop" x1="0%" y1="0%" x2="55%" y2="100%">
      <stop offset="0%" stop-color="#ffffff"/><stop offset="100%" stop-color="#eef2f5"/>
    </linearGradient>
    <linearGradient id="faceLeft" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#dde4e9"/><stop offset="100%" stop-color="#c6d0d7"/>
    </linearGradient>
    <linearGradient id="faceRight" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#bcc7cf"/><stop offset="100%" stop-color="#a1aeb8"/>
    </linearGradient>
    <radialGradient id="badge" cx="34%" cy="26%" r="84%">
      <stop offset="0%" stop-color="#7DEFC0"/><stop offset="46%" stop-color="#10B981"/><stop offset="100%" stop-color="#046b41"/>
    </radialGradient>
    <radialGradient id="badgeHi" cx="50%" cy="50%" r="50%">
      <stop offset="0%" stop-color="#fff" stop-opacity="0.7"/><stop offset="100%" stop-color="#fff" stop-opacity="0"/>
    </radialGradient>
    <clipPath id="sq"><rect width="256" height="256" rx="58"/></clipPath>
    <filter id="padShadow" x="-40%" y="-40%" width="180%" height="180%">
      <feDropShadow dx="0" dy="11" stdDeviation="10" flood-color="#063a30" flood-opacity="0.35"/>
    </filter>
    <filter id="badgeShadow" x="-70%" y="-70%" width="240%" height="240%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#053b2e" flood-opacity="0.45"/>
    </filter>
  </defs>

  <!-- 背景容器（圆角透明角 + 光影） -->
  <g clip-path="url(#sq)">
    <rect width="256" height="256" fill="url(#bg)"/>
    <rect width="256" height="256" fill="url(#botShade)"/>
    <rect width="256" height="256" fill="url(#topGloss)"/>
    <rect x="1.5" y="1.5" width="253" height="253" rx="56.5" fill="none" stroke="#fff" stroke-opacity="0.22" stroke-width="2"/>
  </g>

  <!-- 地面接触阴影 -->
  <ellipse cx="128" cy="218" rx="92" ry="18" fill="#04332b" opacity="0.30"/>

  <!-- 等距记事本（盒体） -->
  <g filter="url(#padShadow)">
    <polygon points="{poly(right)}" fill="url(#faceRight)"/>
    <polygon points="{poly(left)}"  fill="url(#faceLeft)"/>
    <polygon points="{poly(top)}"   fill="url(#faceTop)"/>
  </g>
  <!-- 顶页亮边 -->
  <polygon points="{poly(left_edge)}"  fill="#ffffff" opacity="0.55"/>
  <polygon points="{poly(right_edge)}" fill="#ffffff" opacity="0.28"/>
  <!-- 侧面页堆线 -->
  {stack_svg}

  <!-- 顶面内容 -->
  <polyline points="{poly(chk)}" fill="none" stroke="#10B981" stroke-width="6" stroke-linecap="round" stroke-linejoin="round"/>
  <polygon points="{poly(title)}" fill="#10B981"/>
  <polygon points="{poly(body[0])}" fill="#C4D0D8"/>
  <polygon points="{poly(body[1])}" fill="#C4D0D8"/>
  <polygon points="{poly(body[2])}" fill="#D3DAE1"/>

  <!-- 3D 球形对勾徽标 -->
  <g filter="url(#badgeShadow)"><circle cx="{BX}" cy="{BY}" r="{BR}" fill="url(#badge)"/></g>
  <ellipse cx="{BX-9}" cy="{BY-10}" rx="16" ry="12" fill="url(#badgeHi)"/>
  <path d="M {BX-13} {BY} l 9 9 l 18 -19" fill="none" stroke="#fff" stroke-width="6.5" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
'''

out = os.path.join(os.path.dirname(os.path.abspath(__file__)), "memo-icon.svg")
with open(out, "w") as f:
    f.write(svg)
print("wrote", out)
