## 接口设计

#### ctrl 接口使用http RESTful 架构

api             | method
--              | --
create          | POST
delet           | DELETE
edit            | PUT
list/search/get | GET

#### 接口的数据类型和命名规则,以user为例

api    | struct name
--     | --
create | UserCreate
edit   | UserUpdate
get    | User


