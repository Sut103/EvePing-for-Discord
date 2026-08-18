# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際に Claude Code (claude.ai/code) へガイダンスを提供するものです。

## 言語

ユーザーへの応答は常に日本語で行うこと。これはチャット出力（説明、要約、質問など）すべてに適用され、ユーザーがどの言語で書いたかにかかわらず日本語で応答する。コード、識別子、ファイルの内容についてはこのルールの対象外であり、通常のエンジニアリング規約に従うこと。コミットメッセージの言語・書式については下記「コミット規約」を参照。

CCPM（PRD、エピック、タスクファイル、進捗ファイルなどのローカル成果物）が生成する文章、およびGitHubへの書き込み（Issue本文、PR本文・タイトル、コメントなど）も同様に日本語で記述すること。ただしコミットメッセージは上記の例外のまま、下記「コミット規約」に従い英語とする。この方針はこのプロジェクト固有の指示であり、`.claude/skills/ccpm/` 配下のCCPMスキル自体（汎用テンプレート）は変更しないこと。

## コミット規約

コミットメッセージは英語で記述し、[Conventional Commits](https://www.conventionalcommits.org/)（`<type>[optional scope]: <description>` の形式、例: `feat:`、`fix:`、`docs:`、`refactor:`、`test:`、`chore:`）に従うこと。

## プロジェクトの現状

Epic「eve-ping」（GitHub Issue #4、`epic:eve-ping` ラベル）に基づき、Go + discordgo による実装が完了している。DiscordのScheduled Eventの前日に、「興味あり」を押したユーザーへ個別DMでリマインドを送る常駐Botである。

## 製品概要

姉妹プロジェクト `7DaysPoll-for-Discord`（投票で日程を決めるDiscord Bot）から、「イベント前日リマインド」機能をスコープ外として切り出したスピンオフ企画。

- **やりたいこと**: DiscordのScheduled Eventの前日に、「興味あり」を押したユーザーへ個別DMでリマインドを送る。
- **通知方法**: チャンネルへのメンションではなく個別DM。ギルド／チャンネルの設定を記憶する必要がなく、Botをステートレスに保てる。
- **処理イメージ**: 1日1回のバッチで、Bot導入済み全ギルドを走査 → 各ギルドのScheduled Eventsを取得・フィルタ（開始日=UTC翌日、status=SCHEDULED/ACTIVE） → 各対象イベントの「興味あり」ユーザーをページネーションで全件取得 → 各ユーザーへ個別DM送信（失敗は個別catchしログ記録、継続）。
- **永続化**: 完全ステートレス。DB／KVストアなし、二重送信対策も未実装（PRDで合意済みの設計判断）。

詳細はGitHub Issue #4（Epic本文）と #5〜#12（タスク）を参照。

## プロジェクト管理ワークフロー（ccpm）

このリポジトリには `ccpm` スキルがインストールされている（`.claude/skills/ccpm/`）。これは PRD → エピック → GitHub Issues → 並列エージェント → 実装完了 という仕様駆動の開発ワークフローである。機能の企画・立案、エピックのタスクへの分解、GitHub Issues への同期、Issue への着手、ステータス・スタンドアップの確認、Issue／エピックのクローズなど、ソフトウェア開発ライフサイクルに関わる作業全般で使用する。

`.claude/skills/ccpm/SKILL.md` および `references/conventions.md` からの要点:

- **要件は頭の中ではなくファイルに置く。** すべての機能は PRD として始まり、技術的なエピックとなり、GitHub Issues に分解され、完全なトレーサビリティを保ちながら並列エージェントによって実行される。
- **TDD はライフサイクル全体で強制される**（red → green）。これは独立したフェーズではなく、各工程に組み込まれている:
  - Plan: 各エピックは、タスク分解前に具体的なテスト戦略（Test Strategy）を定義する。
  - Structure: 各タスクは、技術詳細（Technical Details）の前に、受け入れ基準ごとに1つのテストケースを持つテスト計画（Test Plan）を定義する。
  - Execute: 各エージェントは、テスト計画の各項目について、失敗するテスト（red）を書いてから実装（green）を行う — 逆の順序は許されない。
  - **適用範囲はアプリケーションコード（Go実装。`cmd/`・`internal/` 配下など）の変更に限る。** PRD／エピック／タスクファイルの作成・更新、CLAUDE.md・READMEなどのドキュメント、CI/ビルド設定の変更のように、アプリケーションコードを伴わない変更にはこのTDD（red→green）の強制は適用しない。この方針はこのプロジェクト固有の指示であり、`.claude/skills/ccpm/` 配下のCCPMスキル自体（汎用テンプレート）は変更しないこと。
- **スクリプト優先の原則**: 推論を必要としない決定的な読み取り専用のステータス確認については、手動で調べるのではなく `.claude/skills/ccpm/references/scripts/` にある bash スクリプトを実行すること（例: `status.sh`、`standup.sh`、`epic-list.sh`、`epic-show.sh <name>`、`epic-status.sh <name>`、`prd-list.sh`、`prd-status.sh`、`search.sh <query>`、`in-progress.sh`、`next.sh`、`blocked.sh`、`validate.sh`）。LLM による推論は、PRD の作成、並列性の分析、エージェントの起動、更新内容の統合など、本当に必要な作業のために温存すること。
- 5つのフェーズ（Plan、Structure、Sync、Execute、Track）にはそれぞれ `.claude/skills/ccpm/references/` 配下に専用のリファレンスドキュメントがある — 各フェーズで作業する前に該当するドキュメントを読むこと。

## ビルド・テストコマンド

```bash
go build ./...
go vet ./...
go test ./...
```

エントリポイントのバイナリのみをビルドする場合は `go build ./cmd/eveping`。

## アーキテクチャ

- `cmd/eveping/main.go` — エントリポイント。環境変数 `EVEPING_DISCORD_TOKEN` からBotトークンを読み込み（未設定ならエラー終了）、discordgoセッションを開始し、`internal/scheduler` を起動する。
- `internal/discordclient` — discordgoとの結合点をインターフェース化した層（`Client` インターフェース: `Guilds`/`ScheduledEvents`/`EventUsers`/`SendDM`）。本番実装（discordgoラッパー、`discord.go`）とテスト用のインメモリFake（`fake.go`）を提供する。以降のロジックは全てこのインターフェース越しにテストする。
- `internal/batch` — バッチのコアロジック。`FilterTargetEvents`（翌日UTC判定+status絞り込みの純粋関数）、`FetchAllInterestedUsers`（ページネーション取得）、`SendReminderDM`（1件送信+エラーハンドリング）、`RunDailyBatch`（上記を組み合わせ全ギルド・全イベント・全ユーザーを走査し `BatchResult` を集計。個別失敗は分離し処理を継続する）。
- `internal/reminder` — DM本文のフォーマット（`FormatReminder`: イベント名・UTC開始日時・イベントURLを含む）。
- `internal/scheduler` — `time.Ticker` を使い、注入可能な周期でコールバックを呼び続ける常駐ループ（本番は24時間、テストはミリ秒単位を注入）。
- `dockerbuild` — アプリケーションコードは含まない、Docker関連ファイル（`Dockerfile`/`.dockerignore`/`docker-compose.yml`/CIワークフロー）の内容をテキストとして検証する静的アサーションテスト専用パッケージ。Dockerデーモンなしで `go test ./...` の一部として実行できる。
- `Dockerfile` / `.dockerignore` — マルチステージビルド（Go公式イメージでビルド → 軽量な実行イメージ）によるコンテナイメージ定義。Botトークンはイメージに埋め込まず、コンテナ起動時の環境変数（`EVEPING_DISCORD_TOKEN`）としてのみ注入する。

DB・KVストアは使用しない完全ステートレス構成。全ての状態はバッチ実行のたびにDiscord APIから取得する。

新しいソースファイルやパッケージを追加した際は、このセクションを実態に合わせて更新すること。

`README.md`にも「ビルド・テストコマンド」「アーキテクチャ」相当の説明がある（対象読者が異なるため文面は重複してよい）。どちらか一方を更新した際は、もう一方の該当箇所も同じ変更内容に追随させ、両ファイルの記述が食い違ったまま放置されないようにすること。
