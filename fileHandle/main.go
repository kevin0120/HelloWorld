package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

//主要讲述文本处理的相关内容,包括字符串/数字/JSON/XML
//文件操作函数大多在os包里
func main() {
	//创建新目录
	_ = os.Mkdir("./fileHandle/goDir", 0777)
	//创建新的多级目录
	os.MkdirAll("fileHandle/goDir/test1/test2", 0777)
	//删除名称为name的目录,当目录下有文件或者其他目录时会出错
	err := os.Remove("fileHandle/goDir")
	if err != nil {
		fmt.Println(err)
	}
	//根据path删除多级子目录
	os.RemoveAll("fileHandle/goDir")

	//创建空白文件
	//os.NewFile 或者os.Create
	newfile, _ := os.Create("fileHandle/newFile")
	//关闭空白文件
	newfile.Close()
	fileInfo, err1 := os.Stat("fileHandle/newFile")

	if err1 != nil && os.IsNotExist(err1) {
		log.Fatal("文件不存在")
	}
	//查看文件的相关信息
	fmt.Println(fileInfo.Name(), fileInfo.Size(), fileInfo.Mode(), fileInfo.IsDir())
	//重命名与移动  os.Rename(originalPath,newPath)
	//删除文件
	os.Remove("fileHandle/newFile")

	//文件存在则打开,否则新建一个
	/*
		// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
		O_RDONLY int = syscall.O_RDONLY // open the file read-only.
		O_WRONLY int = syscall.O_WRONLY // open the file write-read-only.
		O_RDWR   int = syscall.O_RDWR   // open the file read-write-read.
		 //The remaining values may be or'ed in to control behavior.
		O_APPEND int = syscall.O_APPEND // append data to the file when writing.
		O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
		O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
		O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
		O_TRUNC  int = syscall.O_TRUNC  // if possible, truncate file when opened.
	*/
	file, _ := os.OpenFile("fileHandle/myFile", os.O_WRONLY|os.O_CREATE, 0666)

	file.WriteString("Hello world!!!!\r\n")
	defer file.Close()
	//复制文件
	//创建新文件
	file1, _ := os.Create("fileHandle/myFile1")
	file2, _ := os.Open("fileHandle/myFile")
	//将file写入file1不成功是因为file是只写的
	//byteWritten, _ := io.Copy(file1, file)

	byteWritten, _ := io.Copy(file1, file2)
	fmt.Printf("文件已复制%d bytes", byteWritten)
	defer file1.Close()
	defer file2.Close()
}

/*
func Create(name string) (*File, error) //Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）。如果成功，返回的文件对象可用于I/O；对应的文件描述符具有O_RDWR模式。如果出错，错误底层类型是*PathError。
    func NewFile(fd uintptr, name string) *File //NewFile使用给出的Unix文件描述符和名称创建一个文件。
    func Open(name string) (*File, error) //Open打开一个文件用于读取。如果操作成功，返回的文件对象的方法可用于读取数据；对应的文件描述符具有O_RDONLY模式。如果出错，错误底层类型是*PathError。
    func OpenFile(name string, flag int, perm FileMode) (*File, error) //OpenFile是一个更一般性的文件打开函数，大多数调用者都应用Open或Create代替本函数。它会使用指定的选项（如O_RDONLY等）、指定的模式（如0666等）打开指定名称的文件。如果操作成功，返回的文件对象可用于I/O。如果出错，错误底层类型是*PathError。
    func Pipe() (r *File, w *File, err error) //Pipe返回一对关联的文件对象。从r的读取将返回写入w的数据。本函数会返回两个文件对象和可能的错误。
    func (f *File) Chdir() error //Chdir将当前工作目录修改为f，f必须是一个目录。如果出错，错误底层类型是*PathError。
    func (f *File) Chmod(mode FileMode) error //Chmod修改文件权限。如果出错，错误底层类型是*PathError。
    func (f *File) Chown(uid, gid int) error //修改文件文件用户id和组id
    func (f *File) Close() error  //Close关闭文件f，使文件不能用于读写。它返回可能出现的错误。
    func (f *File) Fd() uintptr //Fd返回与文件f对应的整数类型的Unix文件描述符。
    func (f *File) Name() string //Name方法返回（提供给Open/Create等方法的）文件名称。
    func (f *File) Read(b []byte) (n int, err error) //Read方法从f中读取最多len(b)字节数据并写入b。它返回读取的字节数和可能遇到的任何错误。文件终止标志是读取0个字节且返回值err为io.EOF。
    func (f *File) ReadAt(b []byte, off int64) (n int, err error) //ReadAt从指定的位置（相对于文件开始位置）读取len(b)字节数据并写入b。它返回读取的字节数和可能遇到的任何错误。当n<len(b)时，本方法总是会返回错误；如果是因为到达文件结尾，返回值err会是io.EOF。
    func (f *File) Readdir(n int) ([]FileInfo, error) //Readdir读取目录f的内容，返回一个有n个成员的[]FileInfo，这些FileInfo是被Lstat返回的，采用目录顺序。对本函数的下一次调用会返回上一次调用剩余未读取的内容的信息。
如果n>0，Readdir函数会返回一个最多n个成员的切片。这时，如果Readdir返回一个空切片，它会返回一个非nil的错误说明原因。如果到达了目录f的结尾，返回值err会是io.EOF。
如果n<=0，Readdir函数返回目录中剩余所有文件对象的FileInfo构成的切片。此时，如果Readdir调用成功（读取所有内容直到结尾），它会返回该切片和nil的错误值。如果在到达结尾前遇到错误，会返回之前成功读取的FileInfo构成的切片和该错误。
    func (f *File) Readdirnames(n int) (names []string, err error) //Readdir读取目录f的内容，返回一个有n个成员的[]string，切片成员为目录中文件对象的名字，采用目录顺序。对本函数的下一次调用会返回上一次调用剩余未读取的内容的信息。
如果n>0，Readdir函数会返回一个最多n个成员的切片。这时，如果Readdir返回一个空切片，它会返回一个非nil的错误说明原因。如果到达了目录f的结尾，返回值err会是io.EOF。
如果n<=0，Readdir函数返回目录中剩余所有文件对象的名字构成的切片。此时，如果Readdir调用成功（读取所有内容直到结尾），它会返回该切片和nil的错误值。如果在到达结尾前遇到错误，会返回之前成功读取的名字构成的切片和该错误。
    func (f *File) Seek(offset int64, whence int) (ret int64, err error) //Seek设置下一次读/写的位置。offset为相对偏移量，而whence决定相对位置：0为相对文件开头，1为相对当前位置，2为相对文件结尾。它返回新的偏移量（相对开头）和可能的错误。
    func (f *File) SetDeadline(t time.Time) error // 设置文件读取和写入时间，超时返回错误
    func (f *File) SetReadDeadline(t time.Time) error //设置文件读取时间
    func (f *File) SetWriteDeadline(t time.Time) error // 设置文件写入时间
    func (f *File) Stat() (FileInfo, error) //Stat返回描述文件f的FileInfo类型值。如果出错，错误底层类型是*PathError。
    func (f *File) Sync() error //Sync递交文件的当前内容进行稳定的存储。一般来说，这表示将文件系统的最近写入的数据在内存中的拷贝刷新到硬盘中稳定保存。
    func (f *File) Truncate(size int64) error //Truncate改变文件的大小，它不会改变I/O的当前位置。 如果截断文件，多出的部分就会被丢弃。如果出错，错误底层类型是*PathError。
    func (f *File) Write(b []byte) (n int, err error) //Write向文件中写入len(b)字节数据。它返回写入的字节数和可能遇到的任何错误。如果返回值n!=len(b)，本方法会返回一个非nil的错误。
    func (f *File) WriteAt(b []byte, off int64) (n int, err error) //将len(b)字节写入文件，从字节偏移开始。它返回写入的字节数和错误，写的时候返回一个错误，当n != len（b）
    func (f *File) WriteString(s string) (n int, err error) //WriteString类似Write，参数为字符串。
*/
