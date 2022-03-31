packr2 build  ./dll.go --legacy
空白内容："\"\""
默认生成内容："\&#34;\&#34;"


packr2 build  ./dll.go 生成不用修改的文件
packr2 clean 删除


packr2 build  ./main.go --legacy 可以生成可执行文件


go 包运行顺序  同一个package的所有import ---最末端的文件的全局变量--最末端的init---上层的init