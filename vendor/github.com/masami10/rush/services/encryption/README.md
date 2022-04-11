## 如何开发

1. 需要启用CGO。

2. 关于IDEA下配置，以goland为例。

   可以编辑Run/Debug Config并在Go tool arguments中加入`-tags release`。

## Release编译

1. 正常模式下不进行权限验证, 需要验证则需要加入以下tags
    ```shell
    go build -tags release
    ```
   详情可见`mode.go`
   
2. linux下使用静态链接库编译到可执行文件
   使用`libhasp_linux_x86_64_<vendor-id>.a`
3. windows使用
   将`libhasp_windows_x64_<vendor-id>.lib`重命名为`libhasp_windows_x64_<vendor-id>.a`
   
   修改`hasp_api.h` 之前默认使用dll调用方式
   ```c
   #if defined(__MINGW32__)
   #  define HASP_CALLCONV __declspec(dllimport) __stdcall
   #else
   ```
   替换为
   ```c
   #if defined(__MINGW32__)
   #  define HASP_CALLCONV
   #else
   ```

4. CI的配置详见`.goreleaser.yml`。

## 一些参数

1. 目前默认使用featureID = 2的feature控制rush程序的权限。主要是控制程序是否可以运行或者程序是否过期。
2. 使用LDK中RO内存的payload来决定rush的功能。(格式确定为json)
    ```json
    {
      "controllers":  10
    }
    ```
    目前行为只有controllers控制控制器的数目。可添加xmlProtocol等控制协议的权利。