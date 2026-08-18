---
name: dockerize
status: backlog
created: 2026-08-18T15:35:12Z
updated: 2026-08-18T15:35:12Z
progress: 0%
prd: .claude/prds/dockerize.md
github: (will be set on sync)
---

# Epic: dockerize

## Overview

EvePingをコンテナイメージとして起動できるようにする。既存のアプリケーションロジック（`cmd/eveping`, `internal/*`）には一切手を入れず、マルチステージビルドの `Dockerfile`・`.dockerignore`・開発用 `docker-compose.yml`・README追記のみを追加する、純粋にデプロイ手段を追加するインフラ変更。

## Architecture Decisions

- **ビルドステージ**: `golang:1.24-alpine`（`go.mod` の `go 1.24.7` に対応する最新の1.24系）。依存は `discordgo`/`gorilla/websocket`/`golang.org/x/crypto`/`golang.org/x/sys` のみでCGO不要なため、`CGO_ENABLED=0` で静的バイナリをビルドする。
- **実行ステージ**: `alpine:3.20` 程度の軽量イメージ。Goツールチェインを含めず、ビルド成果物のバイナリのみをコピーする。distrolessも候補だが、初期導入ではシェルが使えて運用時のデバッグが容易なalpineを採用し、distrolessへの移行は必要になれば別タスクとする。
- **非rootユーザー**: 実行ステージ内に専用ユーザーを作成し、そのユーザーでプロセスを起動する。
- **レイヤーキャッシュ**: `go.mod`/`go.sum` を先にCOPYして `go mod download` を実行し、その後ソース全体をCOPYしてビルドすることで、依存変更のない再ビルドを高速化する。
- **秘匿情報の扱い**: `EVEPING_DISCORD_TOKEN` はDockerfile中のARG/ENVに一切埋め込まず、`docker run -e` またはCompose の `environment`/`env_file` を通じてのみコンテナに渡す。
- **ビルドコンテキスト**: `.dockerignore` を追加し、`.git`・ローカル`.env`・`.claude/` 等の不要ファイルをビルドコンテキストから除外する。

## Technical Approach

### Backend Services
既存のGoアプリケーションコード（`cmd/eveping/main.go`, `internal/*`）は変更しない。Docker化はデプロイ手段の追加のみで、アプリケーションの振る舞いに変更はない。

### Infrastructure
- `Dockerfile`（マルチステージ: build → runtime）をリポジトリルートに追加
- `.dockerignore` を追加
- 開発用 `docker-compose.yml` を追加（`.env` からのトークン注入をサポート）
- `README.md` にDockerでのビルド・起動手順を追記

## Test Strategy

### Test Types & Tools
- **静的アサーションテスト（自動・`go test`実行）**: 新設のテストファイルから、リポジトリルートの `Dockerfile` / `docker-compose.yml` をテキストとして読み込み、以下を機械的に検証する。
  - `Dockerfile` に複数の `FROM` 命令があること（マルチステージビルドになっていること）
  - `Dockerfile` 内に `EVEPING_DISCORD_TOKEN` の値がハードコードされていないこと（`ENV EVEPING_DISCORD_TOKEN=<値>` のような記述が存在しないこと）
  - `Dockerfile` に `CGO_ENABLED=0` が設定されていること
  - `Dockerfile` が非rootユーザーで実行される設定（`USER` 命令）を含むこと
  - `docker-compose.yml` が `EVEPING_DISCORD_TOKEN` を環境変数経由（`environment`/`env_file`）で参照しており、値をハードコードしていないこと
  - これらのテストはこの環境（Dockerデーモン不在のサンドボックス）でも `go test ./...` の一部として実行でき、CIでも同様に実行可能。
- **手動スモークテスト（Dockerデーモンが利用可能な環境でのみ実施、README/タスクのDefinition of Doneに手順を明記）**:
  - `docker build .` が成功すること
  - `docker run -e EVEPING_DISCORD_TOKEN=<有効なトークン> <image>` が `go run ./cmd/eveping` と同等に動作し、Discordへ接続してスケジューラが起動すること（既存READMEの手動検証手順に準拠）
  - `EVEPING_DISCORD_TOKEN` を指定せずに起動した場合、`main.go` の既存仕様どおり即座にエラー終了すること
  - `docker compose up`（または `docker compose config`）でCompose定義が問題なく解釈され、`.env` からトークンが注入されること

### Coverage Expectations
- PRDの各Functional Requirementは、上記の自動静的テストまたは手動検証手順のいずれかで必ずカバーする。
- 既存の `go build ./...` / `go vet ./...` / `go test ./...` が新規ファイル追加後も引き続き全て成功すること（回帰がないこと）。

### TDD Notes
- 各タスクは、対応する静的アサーションテストを先に書き（例: 対象ファイルが存在しない/条件を満たさないために失敗する状態を確認 = red）、その後 `Dockerfile`/`docker-compose.yml`/README側を実装してテストをgreenにする。
- Dockerデーモンが利用できないサンドボックス環境では実際のビルド・起動は自動テストできないため、その検証は各タスクの Definition of Done に手動検証ステップとして明記し、自動テストの対象外であることを明示する（既存README「開発用Discordサーバーでの手動検証手順」と同じ考え方）。

## Implementation Strategy

タスク数が少なく依存関係も単純なため、逐次作成する（Small epic, <5 tasks）。Dockerfile本体の確定が他タスクの前提になるため、まずDockerfile/.dockerignoreを完成させ、その後docker-compose.ymlとREADME追記を進める。docker-compose.ymlとREADME追記は互いに独立しており並行可能。

## Task Breakdown Preview

- [ ] Task 1: マルチステージ `Dockerfile` + `.dockerignore` の追加（静的アサーションテスト含む）
- [ ] Task 2: 開発用 `docker-compose.yml` の追加（`.env` からのトークン注入、静的アサーションテスト含む）
- [ ] Task 3: README へのDockerビルド・起動手順の追記

## Dependencies

- Task 2・Task 3 は Task 1（Dockerfileの存在・仕様確定）に依存する。
- 外部依存: Dockerデーモンでのビルド・起動確認（手動検証）にはDocker環境が必要。

## Success Criteria (Technical)

- 新規追加した静的アサーションテストが全てgreenになる。
- `go build ./...` / `go vet ./...` / `go test ./...` が全て成功する。
- （Docker環境がある場合の手動検証）`docker build .` が成功し、`docker run -e EVEPING_DISCORD_TOKEN=...` で `go run ./cmd/eveping` と同等にBotが起動する。
- README の記述に従うだけでDocker起動ができる。

## Estimated Effort

- Size: S
- Hours: 4-6
