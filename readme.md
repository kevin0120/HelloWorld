# go releaser
```bash
go releaser 时tag 必须满足一定的条件 v1.1.1 语义版本号（Semantic Versioning） 1.4.6-beta v是不强制的
https://www.cnblogs.com/walterlv/p/10236470.html

或者 goreleaser release --snapshot 不加标签
https://goreleaser.com/customization/templates/
```


# go get

```bash
go get 出来的tag必须是 v1.1.1格式的tag才可以

go get 会build根目录下的main文件，生成exe放在gopath/bin目录下

go get 后面可以加分支比如github.com/xuri/excelize/v2 取最新的tag

go get 版本管理
https://zhuanlan.zhihu.com/p/105556877

go get github.com/foo

go get -u 不是特别灵敏 最好是@指定tag 这样最准确

https://blog.csdn.net/weixin_28903391/article/details/112100117
https://blog.csdn.net/liuqun0319/article/details/103213396
原因是GOPROXY的存在

# 最新的次要版本或者修订版本(x.y.z, z是修订版本号， y是次要版本号)
go get -u github.com/foo

# 升级到最新的修订版本
go get -u=patch github.com/foo

# 指定版本，若存在tag，则代行使用
go get github.com/foo@v1.2.3

# 指定分支
go get github.com/foo@master

# 指定git提交的hash值
```

#git submodule
```bash
git submodule可以指定分支 ：git submodule add -b v2 git@github.com:qax-os/excelize.git testExcel/fromandtoexcel/source
```
#git lfs
```bash
git lfs git 大文件处理
https://www.jianshu.com/p/493b81544f80
```