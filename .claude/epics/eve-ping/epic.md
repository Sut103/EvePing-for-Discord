---
name: eve-ping
status: backlog
created: 2026-08-17T15:55:13Z
progress: 0%
prd: .claude/prds/eve-ping.md
github: (will be set on sync)
---

# Epic: eve-ping

## Overview

Go + discordgo による常駐Discord Botとして実装する。Botは起動時にDiscordセッションを確立し、内部の `time.Ticker` で1日1回、以下を実行する: 全ギルド走査 → 各ギルドのScheduled Events取得・フィルタ（開始日=UTC翌日、status=SCHEDULED/ACTIVE） → 各対象イベントの興味ありユーザーをページネーションで全件取得 → 各ユーザーへ個別DM送信（失敗は個別catchしログ記録、継続）。DB・KVストアは使わず、全状態はDiscord APIから毎回取得する。

## Architecture Decisions

- **言語/ライブラリ**: Go + `bwmarrin/discordgo`（Discord API操作全般）。7DaysPollと技術スタックを揃え、知見を再利用する。
- **状態管理**: 完全ステートレス。プロセス内メモリにも永続キャッシュを持たない（二重送信対策は実装しない、PRDで合意済み）。
- **スケジューリング**: 外部cronではなく、Bot起動後に `time.Ticker` で内部タイマーを回す常駐プロセス方式。ティッカー周期は24時間、初回実行時刻の考え方は実装タスクで確定する。
- **テスト容易性のための境界分離**: Discord APIとの通信部分（discordgoのセッション呼び出し）をインターフェースの背後に隠し、「どのイベントが対象か」「どのユーザーに送るか」を決めるロジックを純粋関数として抽出する。これによりDiscord APIをモックしたテーブル駆動ユニットテストが書ける。
- **日付比較**: `time.Now().UTC()` を起点に「翌日」を判定。イベントの `ScheduledStartTime` をUTCに変換して年月日で比較する。
- **エラーハンドリング方針**: 1ユーザー・1イベント・1ギルド単位でエラーを分離。goroutineは使わず、シンプルなforループ内でエラーを記録して継続する（v1では並列化しない — スコープはBotが導入される想定ギルド数において十分高速）。

## Technical Approach

### Backend Services（Botプロセス本体、フロントエンドなし）

- **`cmd/eveping/main.go`**: エントリポイント。環境変数からBotトークンを読み込み、discordgoセッションを開始し、`internal/batch` のスケジューラを起動する。
- **`internal/discordclient`**: discordgoセッションをラップするインターフェース（`Guilds() []*Guild`, `ScheduledEvents(guildID) ([]*Event, error)`, `EventUsers(guildID, eventID, after string, limit int) ([]*User, error)`, `SendDM(userID, message string) error`）。本番実装はdiscordgoを呼び出し、テストではフェイク実装に差し替える。
- **`internal/batch`**: バッチのコアロジック。
  - `FilterTargetEvents(events []*Event, now time.Time) []*Event` — 開始日=翌日UTC かつ status が SCHEDULED/ACTIVE のイベントを抽出する純粋関数。
  - `FetchAllInterestedUsers(client, guildID, eventID) ([]*User, error)` — ページネーションを内部でループして全件取得する。
  - `RunDailyBatch(client, now time.Time) BatchResult` — 上記を組み合わせ、全ギルド・全対象イベント・全ユーザーに対してDM送信を試行し、成功/失敗件数を集計して返す。個別エラーは記録するが処理を止めない。
- **`internal/reminder`**: DMメッセージのフォーマット（イベント名・開始日時・イベントURLを含む文面生成）。
- **`internal/scheduler`**: `time.Ticker` を使い24時間ごとに `RunDailyBatch` を呼び出す常駐ループ。

### Infrastructure

- 実行環境: 常駐プロセスとして単一バイナリを稼働させる（VPS等を想定、7DaysPollと同様）。CI/CDやコンテナ化はスコープ外（必要になれば別途）。
- 設定: Botトークンなど機微情報は環境変数で注入（`.env` はローカル開発用、リポジトリにはコミットしない）。
- 外部データストアなし。

## Test Strategy

### Test Types & Tools

- **フレームワーク**: Go標準の `testing` パッケージ + テーブル駆動テスト。追加の依存は最小限にし、必要なら `stretchr/testify` の `assert`/`require` のみ許可する。
- **ユニットテスト**: `internal/batch`（イベントフィルタリングロジック、ページネーション集約ロジック、バッチ集計ロジック）、`internal/reminder`（メッセージフォーマット）を対象に、`internal/discordclient` のフェイク実装を注入して外部API呼び出しなしで検証する。
- **統合テスト**: discordgo実クライアントを使った自動テストは対象外（実Discordサーバーが必要なため）。代わりに `internal/discordclient` のインターフェース実装が正しくdiscordgoの型・エラーをマッピングしているかを、責務の狭い単体テストで担保する。
- **手動検証**: 実装完了時に開発用Discordサーバーでテスト用イベントを作成し、実際にDMが届くことを目視確認する（自動テストの範囲外として明記）。

### Coverage Expectations

- `internal/batch` のフィルタリング・集計ロジックは分岐網羅（日付境界: 前日/当日/翌日/翌々日、status: SCHEDULED/ACTIVE/CANCELED/COMPLETED、ページネーション: 0件/1ページ/複数ページ、送信成功/個別失敗が他に波及しないこと）を必須とする。
- `internal/reminder` のメッセージ生成は正常系1件以上。
- カバレッジ数値目標は設けないが、Test Plan（各タスクの受け入れ基準ごとに1テストケース）を満たすことを必須とする。

### TDD Notes

各タスクはTest Planに列挙したテストケースについて、red（失敗確認）→green（最小実装）の順で進める。Discord APIとの結合部分は `internal/discordclient` インターフェースのフェイクを先に用意し、それを使ってバッチロジックのテストを先に書けるようにする。

## Implementation Strategy

- タスクは依存関係の少ない順に、まずプロジェクト雛形とdiscordクライアント抽象化を作り、その後バッチロジック（フィルタ→ページネーション→送信→集計）を積み上げ、最後にスケジューラとエントリポイントで結線する。
- リスク軽減: Discord APIとの結合を最初にインターフェース化することで、以降のロジックタスクはモックで並列開発できるようにする。
- 段階的デリバリー: 各タスクは単体でビルド・テストが通る状態を保ちながら積み上げる（大きな一枚岩PRを避ける）。

## Task Breakdown Preview

- [ ] **Task 1: プロジェクト雛形 + discordclientインターフェース定義** — go.mod、ディレクトリ構成、discordgoセッション初期化、`internal/discordclient` インターフェースとdiscordgo実装・テスト用フェイクの雛形
- [ ] **Task 2: イベントフィルタリングロジック** — `FilterTargetEvents`（翌日UTC判定 + status絞り込み）とそのユニットテスト
- [ ] **Task 3: 興味ありユーザーのページネーション取得** — `FetchAllInterestedUsers` とそのユニットテスト（0件/1ページ/複数ページ）
- [ ] **Task 4: リマインドメッセージ生成** — `internal/reminder` のフォーマットロジックとテスト
- [ ] **Task 5: DM送信 + 個別エラーハンドリング** — 送信失敗を捕捉し処理を継続するロジックとテスト（エラーコード50007等を含むケース）
- [ ] **Task 6: バッチ実行の統合（RunDailyBatch）** — Task2〜5を組み合わせ、全ギルド・全イベント・全ユーザーを走査して結果を集計するロジックとテスト
- [ ] **Task 7: スケジューラ + エントリポイント結線** — `time.Ticker` による24時間周期実行、環境変数からのトークン読み込み、ログ出力、main.goの結線
- [ ] **Task 8: README・動作確認手順の整備** — セットアップ手順、必要なBot権限/Intent、手動検証手順のドキュメント化

## Dependencies

- `github.com/bwmarrin/discordgo`（外部ライブラリ）
- Discord Developer Portalでのアプリケーション作成・Botトークン発行・`GUILD_SCHEDULED_EVENTS` を含む必要Intentの有効化（人手による準備、実装タスクの前提）
- Go 1.21以降のツールチェイン

## Success Criteria (Technical)

- `go build ./...` が成功し、`go vet` / `go test ./...` がクリーンに通る
- `internal/batch` の分岐網羅テスト（日付境界・status・ページネーション・個別失敗の非波及）が全てgreen
- 開発用Discordサーバーでの手動検証で、翌日開始イベントの興味ありユーザーにDMが届くことを確認済み
- 外部DB・KVストアへの依存が最終実装に一切存在しない（コードレビューで確認）

## Estimated Effort

- 全体で1〜2人日程度を想定（Task 1〜8、それぞれ数時間規模）。
- クリティカルパス: Task 1（雛形）→ Task 2/3/4/5（並列可能）→ Task 6（統合）→ Task 7（結線）→ Task 8（ドキュメント）。

## Tasks Created
- [ ] 001.md - プロジェクト雛形 + discordclientインターフェース定義（parallel: true）
- [ ] 002.md - イベントフィルタリングロジック（parallel: true, depends_on: 1）
- [ ] 003.md - 興味ありユーザーのページネーション取得（parallel: true, depends_on: 1）
- [ ] 004.md - リマインドメッセージ生成（parallel: true, depends_on: 1）
- [ ] 005.md - DM送信 + 個別エラーハンドリング（parallel: true, depends_on: 1）
- [ ] 006.md - バッチ実行の統合（RunDailyBatch）（parallel: false, depends_on: 2,3,4,5）
- [ ] 007.md - スケジューラ + エントリポイント結線（parallel: false, depends_on: 6）
- [ ] 008.md - README・動作確認手順の整備（parallel: false, depends_on: 7）

Total tasks: 8
Parallel tasks: 5（001〜005、うち001が他4件の前提）
Sequential tasks: 3（006, 007, 008）
Estimated total effort: 15 hours
