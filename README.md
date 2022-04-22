# pinyin-search
<a href="https://github.com/jeessy2/pinyin-search/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/release/jeessy2/pinyin-search.svg?logo=github&style=flat-square"></a> <img src=https://goreportcard.com/badge/github.com/jeessy2/pinyin-search /> <img src=https://img.shields.io/docker/image-size/jeessy/pinyin-search /> <img src=https://img.shields.io/docker/pulls/jeessy/pinyin-search />

提供一些接口, 支持一段中英文进行分词, 分词后的数据转换为拼音, 并保存到search中, 目前只实现了meilisearch.


## 安装
```
docker run -d \
  --name pinyin-search \
  --restart=always -p 7701:7701 \
  jeessy/pinyin-search \
  -meiliHost https://meiliHost.com -meiliKey xxxx \
  -auth xxxx
```

## 支持接口验证

http接口中的header中需要传`Authorization`, 参数为安装时指定的auth

## 添加/更新数据接口(更新通过dataId删除之前的索引)

```
http://localhost:7701/addUpdate?tenant=projectName&indexName=test&dataId=123456789&data=今天天气真好啊
```

## 删除接口 (通过dataId可删除索引)

```
http://localhost:7701/delete?tenant=projectName&indexName=test&dataId=123456789
```

## 搜索建议接口

```
http://localhost:7701/suggest?tenant=projectName&indexName=test&data=jttq
```