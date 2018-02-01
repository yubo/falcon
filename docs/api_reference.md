### falcon API Reference


This is a generated documentation. Please read the proto files for more.


##### message `DataPoint` (tsdb/tsdb.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| key |  | Key |
| value |  | TimeValuePair |



##### message `DataPoints` (tsdb/tsdb.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| key |  | Key |
| values |  | (slice of) TimeValuePair |



##### message `Key` (tsdb/tsdb.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| key |  | bytes |
| shardId |  | int32 |



##### message `TimeValuePair` (tsdb/tsdb.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| timestamp |  | int64 |
| value |  | double |



##### service `Agent` (agent/agent.proto)

| Method | Request Type | Response Type | Description |
| ------ | ------------ | ------------- | ----------- |
| Put | PutRequest | PutResponse |  |
| GetStats | Empty | Stats |  |
| GetStatsName | Empty | StatsName |  |



##### message `Empty` (agent/agent.proto)

Empty field.



##### message `Item` (agent/agent.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| metric |  | bytes |
| tags |  | bytes |
| type |  | bytes |
| value |  | double |
| timestamp |  | int64 |



##### message `PutRequest` (agent/agent.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| items |  | (slice of) Item |



##### message `PutResponse` (agent/agent.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| n |  | int32 |



##### message `Stats` (agent/agent.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counter |  | (slice of) uint64 |



##### message `StatsName` (agent/agent.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counterName |  | (slice of) bytes |



##### service `Alarm` (alarm/alarm.proto)

| Method | Request Type | Response Type | Description |
| ------ | ------------ | ------------- | ----------- |
| Put | PutRequest | PutResponse |  |
| GetStats | Empty | Stats |  |
| GetStatsName | Empty | StatsName |  |



##### message `Empty` (alarm/alarm.proto)

Empty field.



##### message `Event` (alarm/alarm.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| tagId |  | int64 |
| key |  | bytes |
| expr |  | bytes |
| msg |  | bytes |
| timestamp |  | int64 |
| value |  | double |
| Priority |  | int32 |



##### message `PutRequest` (alarm/alarm.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| events |  | (slice of) Event |



##### message `PutResponse` (alarm/alarm.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| n |  | int32 |



##### message `Stats` (alarm/alarm.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counter |  | (slice of) uint64 |



##### message `StatsName` (alarm/alarm.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counterName |  | (slice of) bytes |



##### service `Service` (service/service.proto)

| Method | Request Type | Response Type | Description |
| ------ | ------------ | ------------- | ----------- |
| Put | PutRequest | PutResponse |  |
| Get | GetRequest | GetResponse |  |
| GetStats | Empty | Stats |  |
| GetStatsName | Empty | StatsName |  |



##### message `Empty` (service/service.proto)

Empty field.



##### message `GetRequest` (service/service.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| keys |  | (slice of) tsdb.Key |
| start |  | int64 |
| end |  | int64 |
| consolFun |  | Cf |



##### message `GetResponse` (service/service.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| data |  | (slice of) tsdb.DataPoints |



##### message `PutRequest` (service/service.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| data |  | (slice of) tsdb.DataPoint |



##### message `PutResponse` (service/service.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| n |  | int32 |



##### message `Stats` (service/service.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counter |  | (slice of) uint64 |



##### message `StatsName` (service/service.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counterName |  | (slice of) bytes |



##### service `Transfer` (transfer/transfer.proto)

| Method | Request Type | Response Type | Description |
| ------ | ------------ | ------------- | ----------- |
| Put | PutRequest | PutResponse |  |
| Get | GetRequest | GetResponse |  |
| GetStats | Empty | Stats |  |
| GetStatsName | Empty | StatsName |  |



##### message `DataPoint` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| key |  | bytes |
| value |  | tsdb.TimeValuePair |



##### message `DataPoints` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| key |  | bytes |
| values |  | (slice of) tsdb.TimeValuePair |



##### message `Empty` (transfer/transfer.proto)

Empty field.



##### message `GetRequest` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| keys |  | (slice of) bytes |
| start |  | int64 |
| end |  | int64 |
| consolFun |  | service.Cf |



##### message `GetResponse` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| data |  | (slice of) DataPoints |



##### message `PutRequest` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| data |  | (slice of) DataPoint |



##### message `PutResponse` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| n |  | int32 |



##### message `Stats` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counter |  | (slice of) uint64 |



##### message `StatsName` (transfer/transfer.proto)

| Field | Description | Type |
| ----- | ----------- | ---- |
| counterName |  | (slice of) bytes |



