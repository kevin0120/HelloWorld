cmake_minimum_required(VERSION 3.21)
project(Cmake)

set(CMAKE_CXX_STANDARD 14)

set(SOURCE_FILES src/library.cpp)

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



add_executable(Cmake main.cpp)
target_link_directories(Cmake PRIVATE ${CMAKE_BINARY_DIR})

#静态编译
target_link_libraries(Cmake PRIVATE spc_libStatic)
#动态编译
#target_link_libraries(Cmake PRIVATE spc_lib)


target_include_directories(Cmake
        PRIVATE
        ./include)