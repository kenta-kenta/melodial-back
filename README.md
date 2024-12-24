# 日記アプリ × 音楽生成

## 概要

日記を書くとその内容に基づいた音楽を生成するアプリケーション．

## 機能一覧

- ユーザー認証
- 日記の生成・編集・削除
- 日記の一覧表示
- 日記に基づく音楽生成
- 生成された音楽の再生・保存

## 技術スタック

### フロントエンド

- React 18
- TypeScript
- TailwindCSS
- Tanstack Query
- Zustand

### バックエンド

- Go (Echo)
- postgreSQL
- JWT

## 環境構築

```bash
# frontend
cd frontend
npm install
npm run dev

# backend
cd backend
go mod tidy
go run main.go

# database
docker compose up -d
```
