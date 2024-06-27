# GO OpenAI API Proxy

このプロジェクトは、Ginフレームワークを使用してOpenAIのAPIにリクエストをプロキシするシンプルなGoアプリケーションです。

## 必要条件

- Go 1.16以上
- OpenAI APIキー
- .envファイル

## インストール

1. このリポジトリをクローンします。

    ```bash
    git clone <リポジトリURL>
    cd <リポジトリディレクトリ>
    ```

2. 必要なGoモジュールをインストールします。

    ```bash
    go mod tidy
    ```

3. `.env`ファイルをプロジェクトのルートディレクトリに作成し、OpenAI APIキーを設定します。

    ```env
    OPENAI_API_KEY=your_openai_api_key
    ```

## 実行

以下のコマンドでアプリケーションを起動します。

    ```bash
    go run main.go
    ```

アプリケーションはデフォルトで`localhost:3000`で起動します。

## エンドポイント

### POST `/api/openai`

このエンドポイントは、受信したJSONリクエストをOpenAIのAPIにプロキシし、レスポンスを返します。

#### リクエスト例

    ```json
    {
        "model": "text-davinci-003",
        "prompt": "Say this is a test",
        "max_tokens": 7
    }
    ```

#### レスポンス例

    ```json
    {
        "id": "cmpl-1234567890",
        "object": "text_completion",
        "created": 1612303456,
        "model": "text-davinci-003",
        "choices": [
            {
                "text": "This is a test.",
                "index": 0,
                "logprobs": null,
                "finish_reason": "length"
            }
        ],
        "usage": {
            "prompt_tokens": 5,
            "completion_tokens": 7,
            "total_tokens": 12
        }
    }
    ```

## 注意点

- このアプリケーションは開発およびテスト目的でのみ使用してください。
- APIキーの漏洩に注意してください。

## ライセンス

このプロジェクトはMITライセンスの下で提供されています。詳細はLICENSEファイルをご覧ください。