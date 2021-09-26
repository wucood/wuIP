# wuIP

IP地址查询（中国IP调用高德个人免费接口，中国以外IP使用GeoLite2）

### 返回JSON数据

```json
# 检测到中国IP,返回以下信息
{"province":"吉林省","city":"长春市","district":"南关区","isp":"中国教育网","ip":"59.72.94.129"}
# 检测到非中国IP,返回以下信息
{"ip":"89.72.94.121","country":"波兰","province":"West Pomerania","city":"什切青","location":["14.547500","53.428600"]}
```
