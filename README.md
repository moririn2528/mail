# mail
AWS SES 使いたくなった

## 方法
- ドメインをお名前ドットコムなどで取得
- route 53 でホストゾーンの作成、NS をドメイン取得サイトに貼る
- SES で create identify, domain を書いて、他はそのままで作成
- CNAME は自動で route 53 に設定されるので verify まで待つ
- 受信側のメールアドレスの create identify も行う。メールアドレスに認証メールが届くので、リンクを踏む。
- IAM ユーザー作成、SES 権限を渡す
- 後は適当にコードを書いてアップロード

## 受信
- route 53 で MX 10 inbound-smtp.us-west-2.amazonaws.com を入れる
- S3 作成
- rule set 作成、active にして、rule 作成、receipt に受信用メールアドレス入れる。action で S3 に投げるを選択。