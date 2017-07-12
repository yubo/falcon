## falcon acl

- global-r `在任何节点拥有 falcon-read token`
- global-w `在任何节点拥有 falcon-opterate token`
- global-a `在任何节点拥有 falcon-admin token`
- tag-r `在当前节点或上级节点拥有 falcon-read token`
- tag-w `在当前节点或上级节点拥有 falcon-opterate token`
- tag-a `在当前节点或上级节点拥有 falcon-admin token`
- owner `自己创建的对象`
- 获取用户global权限接口 `http://dev02:8001/v1.0/auth/info`
  *  "admin": global-a
  *  "operator": global-w
  *  "reader":  global-r

-----

|module      | add/edit/del                        | list/get/search   |
|------------|-------------------------------------|-------------------|
| tag        | tag-w(parent)                       |  tag-r            |
| tag-tpl    | tag-w                               |  global-r         |
| tag-host   | tag-w                               |  global-r         |
| aggreator  | tag-w                               |  global-r         |
| plugin     | tag-w                               |  global-r         |
| alarm      | owner(edit/del) / global-w(add)     |  global-r         |
| template   | owner(edit/del) / global-w(add)     |  global-r         |
| expression | owner(edit/del) / global-w(add)     |  global-r         |
| nodata     | owner(edit/del) / global-w(add)     |  global-r         |
| team       | owner(edit/del) / global-w(add)     |  global-r         |
| dashboard  | global-r                            |  global-r         |
| admin, rel, token, role, user|  tag-a            |  global-a         |
