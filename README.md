# Color Game

一個互動式的顏色配對遊戲。

## 🎮 在線 Demo

[點擊這裡體驗遊戲](http://13.231.89.103)

## 🌟 特色

- 即時顏色配對
- WebSocket 通訊
- 響應式設計
- Docker 容器化部署

## 🛠 技術棧

### 前端
- Vue 3 (Composition API)
- WebSocket
- Vite
- Tailwind CSS

### 後端
- Go
- Gorilla WebSocket
- Docker

### 部署
- AWS EC2
- Docker Compose
- GitHub Actions CI/CD
- Nginx

## 🔄 CI/CD

使用 GitHub Actions 實現自動化部署：
- 自動構建 Docker 映像
- 自動部署到 AWS EC2
- 自動版本追踪（Git commit hash）

## 🎯 項目亮點

1. **即時互動**
   - 使用 WebSocket 實現即時反饋
   - 自定義的斷線重連機制

2. **可擴展架構**
   - 前後端分離
   - 容器化部署
   - 模塊化設計

3. **自動化部署**
   - 完整的 CI/CD 流程
   - 零停機部署
   - 自動版本追踪

4. **用戶體驗**
   - 響應式設計
   - 直觀的遊戲界面
   - 即時反饋

## 📝 開發經驗總結

1. **技術選型**
   - Vue 3 的 Composition API 提供更好的代碼組織
   - Go 的高性能特性適合 WebSocket 服務
   - Docker 簡化部署和環境一致性

2. **架構設計**
   - 前後端分離提高可維護性
   - WebSocket 保證實時性
   - 容器化便於擴展

3. **部署優化**
   - 使用 GitHub Actions 自動化部署
   - Docker Compose 簡化服務編排
   - Nginx 處理靜態資源和反向代理

## 🚀 未來優化方向

1. **功能擴展**
   - 添加多人遊戲模式
   - 添加排行榜系統

2. **技術優化**
   - 引入 Redis 緩存

3. **用戶體驗**
   - 優化移動端適配

