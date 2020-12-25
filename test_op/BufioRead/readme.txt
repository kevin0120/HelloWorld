bufio.Scanner

A: 最初分配一个4KB的切片用于存储读到的数据

B：用byte[0:4096]去read 要么阻塞，要摸极端情况读出100个空的信号，要么正常读取n>0,end=n

		极端情况：设置err 下次循环的时候退出

		n<4096 且里面有'\n':  start=n， token=byte[0:n]中的处理结果 返回token   继续循环用byte[n:4096]去read tocken=byte[start:end]中的处理结果。
					当byte满了或者start>2048时，将[start:end]copy dao[0:end-start] 重置start和end 绝大多数是这种机制


		n=4096 且里面没有'\n':先把byte[0:start]的空间copy掉 去read一次 若byte[0:4096]还是没有'\n'，则扩容 2倍的扩容，直到>64KB则done下此start=end=0

