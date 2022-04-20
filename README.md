# pinyin-search
<a href="https://github.com/jeessy2/pinyin-search/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/release/jeessy2/pinyin-search.svg?logo=github&style=flat-square"></a> <img src=https://goreportcard.com/badge/github.com/jeessy2/pinyin-search /> <img src=https://img.shields.io/docker/image-size/jeessy/pinyin-search /> <img src=https://img.shields.io/docker/pulls/jeessy/pinyin-search />

提供一些接口, 支持一段中英文进行分词, 分词后的数据转换为拼音, 并保存到search中, 目前只实现了meilisearch.

## 添加/更新数据接口(通过guid可删除以前的索引)

```
http://localhost:7701/addUpdate?tenant=test&indexName=test&guid=123456789&data=今天天气真好啊
```

## 删除接口

```
http://localhost:7701/delete?tenant=test&indexName=test&guid=123456789
```

## 搜索建议接口

```
http://localhost:7701/suggestion?data=jttq&indexName=test
```