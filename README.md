# cli-oauth
cliでoauthを試す（Twitter）

**これはあくまでサンプルです**

## Usage

```sh
$ go run main.go login
```

で、http://localhost:3000 に認証するためのローカルサーバーが立ち上がります。
認証を終えると、`.token.json` というファイルが生成されていて、そこに`AT`と`ST`が保存されています。
あとは`CK`と`CS`の4つのキーで基本認証するのでAPIが自由に扱えます。

試す場合

action/cli.goの`Tweet`関数の中の`twitter.Tweet`の引数を自由に書き換えてから

```sh
$ go run main.go tweet
```

を行ってください。