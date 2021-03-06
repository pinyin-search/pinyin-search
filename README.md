# pinyin-search
<a href="https://github.com/pinyin-search/pinyin-search/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/release/pinyin-search/pinyin-search.svg?logo=github&style=flat-square"></a> <img src=https://goreportcard.com/badge/github.com/pinyin-search/pinyin-search /> <img src=https://img.shields.io/docker/image-size/pinyinsearch/pinyin-search /> <img src=https://img.shields.io/docker/pulls/pinyinsearch/pinyin-search />

提供一些接口, 支持一段中英文进行分词, 分词后的数据转换为拼音, 并保存到search中, 目前只实现了meilisearch.


## 安装
```
docker run -d \
  --name pinyin-search \
  --restart=always -p 7701:7701 \
  pinyinsearch/pinyin-search \
  -meiliHost https://meiliHost.com -meiliKey xxxx \
  -auth xxxx
```

## 支持接口验证

http接口中的header中需要传`Authorization`, 参数为安装时指定的auth

## 更新接口(不存在会新增, 更新会通过dataId删除之前的索引)

```
http://localhost:7701/update?tenant=projectName&indexName=test&dataId=123456789&data=今天天气真好啊
```

## 批量更新接口
```
http://localhost:7701/updateBatch

Content-Type: application/json
```

``` json
[
  {"tenant": "projectName", "indexName": "test", "dataId":  "1", "data":  "我是帅哥"},
  {"tenant": "projectName", "indexName": "test", "dataId":  "2", "data":  "我是美女"}
]
```

## 删除接口 (通过dataId可删除索引)

```
http://localhost:7701/delete?tenant=projectName&indexName=test&dataId=123456789
```

## 删除所有

```
http://localhost:7701/deleteAll?tenant=projectName&indexName=test
```

## 搜索建议接口

```
http://localhost:7701/suggest?tenant=projectName&indexName=test&data=jttq
```