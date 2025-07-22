# tx-disguise

一個終端機工具，能在顯示台灣期貨交易所（TAIFEX）期貨與現貨即時價格的同時，偽裝交易行為，並展示假系統資訊。以 Go 與 Bubble Tea TUI 框架打造。

## 功能特色
- **終端介面**：即時顯示系統進程資訊（來自 `top`）與期貨/現貨價格。
- **期貨支援**：支援小台（MXF）、微台（TMF）及預設 TXF 合約。
- **可設定**：使用簡易 YAML 設定檔（目前預留，未啟用）。
- **程式碼檢查**：整合 golangci-lint，確保程式品質。
- **自動發佈**：GitHub Actions 自動建置並發佈 Darwin/arm64（Apple Silicon）二進位檔。

## 使用說明

```
用法: txd [-v] [-h] [ -y | -z ]
    -v: 顯示版本
    -h: 顯示說明
合約選項:
    -y: 小台 (MXF)
    -z: 微台 (TMF)
範例:
    ./txd -y
```

- 下載對應平台的發佈檔案（如 `tx-disguise-vX.Y.Z-arm64.tar.gz`）。
- 解壓縮：
  ```sh
  tar -xzvf tx-disguise-vX.Y.Z-arm64.tar.gz
  ```
- 執行程式：
  ```sh
  ./txd [flags]
  ```

### 終端介面操作

執行 `txd` 後，會看到兩個主要畫面：

- **系統資訊畫面**：顯示假系統進程資訊（來自 `top`）及目前期貨/現貨價格。
- **歷史紀錄畫面**：可捲動瀏覽過去一小時的期貨/現貨價格。

**操作鍵：**
- `[q]` 離開程式
- `[tab]` 切換系統資訊與歷史紀錄畫面
- `[up/down]` 在歷史紀錄中捲動

系統資訊每 0.5 秒更新一次，期貨價格每 2 秒更新一次，並每分鐘加入一次歷史紀錄。

## 開發

### 環境需求
- Go 1.24+

### 開發模式執行

```
make dev-tui
```
直接以 `go run ./cmd/tx-disguise` 執行 TUI。

### 程式碼檢查

```
make lint
```
自動安裝 `golangci-lint` 並執行格式化與檢查。

## 專案結構
- `cmd/tx-disguise/main.go` — 進入點
- `internal/futures/` — 期貨邏輯、API 與工具
- `internal/tui/` — 終端介面（Bubble Tea）
- `config/config.yaml` — 設定檔（預留）
- `Makefile` — 開發與檢查任務
- `.air.toml` — Air 熱重載設定
- `.github/workflows/release.yml` — 自動發佈

## 授權
MIT
