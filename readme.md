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
git rm  --cached path_to_submodule
git submodule update --init --recursive

git submodule init
# Submodule 'xxx/xxx' (http://xxxx/xxx/xxx.git) registered for path 'xxx/xxx'
git submodule sync
# Synchronizing submodule url for 'xxx/xxx'
git submodule update
#Cloning into '/Users/xxx/xxx/xxx'...
#Submodule path 'xxx/xxx': checked out '39cabde3d5c8aeba5623424asd7f5948e7f515f9f28db'

每个submodule 都运行
git submodule foreach git checkout .

```
#git lfs
```bash
git lfs git 大文件处理
https://www.jianshu.com/p/493b81544f80

执行 git lfs install 开启lfs功能
使用 git lfs track 命令进行大文件追踪 例如git lfs track "*.png" 追踪所有后缀为png的文件
使用 git lfs track 查看现有的文件追踪模式
提交代码需要将gitattributes文件提交至仓库. 它保存了文件的追踪记录
提交后运行git lfs ls-files 可以显示当前跟踪的文件列表
将代码 push 到远程仓库后，LFS 跟踪的文件会以『Git LFS』的形式显示:
clone 时 使用'git clone' 或 git lfs clone均可

```

#to do
```bash
fsnotify
nats
modbus
```