## falcon ctrl auth

#### google
- https://developers.google.com/identity/protocols/OpenIDConnect#getcredentials
- https://console.developers.google.com/apis/credentials/oauthclient

- 网页应用的客户端id
- 已获授权的 javascript 来源
  - http://falcon.pt.xiaomi.com
  - https://falcon.pt.xiaomi.com
- 已获授权的重定向 URI
  - http://falcon.pt.xiaomi.com/v1.0/auth/callback/google
  - https://falcon.pt.xiaomi.com/v1.0/auth/callback/google

将得到的客户端id 和 客户端密钥 分别填写到
falcon.conf -> ctrl -> GoogleClientID , GoogleClientSecret



#### github
- https://help.github.com/enterprise/2.11/admin/guides/user-management/using-github-oauth/
- https://github.com/settings/applications/


将得到的客户端id 和 客户端密钥 分别填写到
falcon.conf -> ctrl -> GithubClientID , GithubClientSecret


#### misso
- 使用内网域名作为 misso 回调地址
