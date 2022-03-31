cmake_minimum_required(VERSION 3.21)
project(Cmake)
#C++测试出来的库文件cgo调用不了,且源文件后缀名要为.c  c++需要与extern "C" 结合使用
#set(CMAKE_CXX_STANDARD 14)
set(CMAKE_C_STANDARD 99)

set(SOURCE_FILES src/library.c)

add_library(spc_lib_obj OBJECT "")
target_compile_options(spc_lib_obj PRIVATE -fPIC)
target_sources(spc_lib_obj
        PRIVATE
        ${SOURCE_FILES}
        )

target_include_directories(spc_lib_obj
        PRIVATE
        include
        )
#动态库文件
add_library(spc_lib SHARED $<TARGET_OBJECTS:spc_lib_obj>)

#静态库文件
add_library(spc_libStatic STATIC $<TARGET_OBJECTS:spc_lib_obj>)
set_target_properties(spc_libStatic PROPERTIES OUTPUT_NAME spc_lib)


#后缀名要对应起来 c或者cpp
add_executable(Cmake main.c)
target_link_directories(Cmake PRIVATE ${CMAKE_BINARY_DIR})

#静态编译
target_link_libraries(Cmake PRIVATE spc_libStatic)
#动态编译
#target_link_libraries(Cmake PRIVATE spc_lib)


target_include_directories(Cmake
        PRIVATE
        ./include)