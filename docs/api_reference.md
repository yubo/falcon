### falcon API Reference


This is a generated documentation. Please read the proto files for more.


##### message `File` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| name |  | string |
| data |  | bytes |



##### message `GetRequest` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| start |  | int64 |
| end |  | int64 |
| host |  | bytes |
| name |  | bytes |
| step |  | int32 |
| consolFun |  | Cf |
| type |  | ItemType |



##### message `GetResponse` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| host |  | bytes |
| name |  | bytes |
| type |  | ItemType |
| step |  | int32 |
| vs |  | (slice of) RRDData |



##### message `GetRrdRequest` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| key |  | bytes |



##### message `GetRrdResponse` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| file |  | File |



##### message `Item` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| value |  | double |
| ts |  | int64 |
| step |  | int32 |
| type |  | ItemType |
| host |  | bytes |
| name |  | bytes |
| tags |  | bytes |



##### message `RRDData` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| ts |  | int64 |
| v |  | double |



##### message `Response` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| code |  | int32 |



##### message `UpdateRequest` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| items |  | (slice of) Item |



##### message `UpdateResponse` (falcon.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| total |  | int32 |
| errors |  | int32 |



##### service `Transfer` (transfer/transfer.proto)

| Method | Request Type | Response Type | Description |
| ------ | ------------ | ------------- | ----------- |
| Update | falcon.UpdateRequest | falcon.UpdateResponse |  |
| Get | falcon.GetRequest | falcon.GetResponse |  |
| GetRrd | falcon.GetRrdRequest | falcon.GetRrdResponse |  |



##### service `Backend` (backend/backend.proto)

| Method | Request Type | Response Type | Description |
| ------ | ------------ | ------------- | ----------- |
| Update | falcon.UpdateRequest | falcon.UpdateResponse |  |
| Get | falcon.GetRequest | falcon.GetResponse |  |
| GetRrd | falcon.GetRrdRequest | falcon.GetRrdResponse |  |



