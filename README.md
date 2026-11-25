# PHPerKaigi 2025 Albatross


# これは何？

2025-03-21 から 2025-03-23 にかけて開催された [PHPerKaigi 2025](https://phperkaigi.jp/2025/) の中の企画、PHPer コードバトルのシステムです。

現在は回答を締め切っていますが、コードゴルフの問題や開催中の参加者の方々の回答は、以下のページから確認できます。

https://t.nil.ninja/phperkaigi/2025/code-battle/

> [!NOTE]
> 現在はアーカイブされた静的ファイルがサーブされています: https://github.com/nsfisis/phperkaigi-2025-albatross-archive


# サンドボックス化の仕組み

ユーザから任意のコードを受け付ける関係上、何も対策をしないと深刻な脆弱性を抱えてしまいます。

このシステムでは、送信されたコードを WebAssembly へ変換された PHP 処理系で実行することで、サンドボックス化を実現しています。


# License

The contents of the repository are licensed under The MIT License, except for

* [backend/admin/assets/css/normalize.css](backend/admin/assets/normalize.css),
* [backend/admin/assets/css/sakura.css](backend/admin/assets/sakura.css),
* [frontend/public/favicon.svg](frontend/public/favicon.svg) and
* [frontend/public/logo.svg](frontend/public/logo.svg).

See [LICENSE](./LICENSE) for copylight notice.
