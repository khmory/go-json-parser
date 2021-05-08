# Go JSON Parser
学習用のGo製JSONパーサ。
実装はshin1x1様の[shin1x1/php8-toy-json-parser](https://github.com/shin1x1/php8-toy-json-parser)を参考にさせていただきました。


## 実行サンプル

```
cd {project}/src
echo '{"hoge":"fuga", "foo":[100, true, null]}' | go run main.go
```