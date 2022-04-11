cmake_minimum_required(VERSION 3.17)
project(spc)
set(CMAKE_CXX_STANDARD 14)

set(SOURCE_FILES src/spc.cpp src/stats.cpp src/type.cpp)

add_library(spc_lib_obj OBJECT "")
target_compile_options(spc_lib_obj PRIVATE -fPIC)
target_sources(spc_lib_obj
        PRIVATE
        ${SOURCE_FILES}
        )

target_include_directories(spc_lib_obj
        PRIVATE
        include
        numcpp/include
        )

add_subdirectory(numcpp)

set(Boost_USE_STATIC_LIBS ON)  # only find static libs
set(Boost_USE_DEBUG_LIBS OFF)  # ignore debug libs and
set(Boost_USE_RELEASE_LIBS ON)  # only find release libs
set(Boost_USE_MULTITHREADED ON)
set(Boost_USE_STATIC_RUNTIME OFF)
find_package(Boost 1.66.0)
if (Boost_FOUND)
    target_include_directories(spc_lib_obj
            PRIVATE
            ${Boost_INCLUDE_DIRS})
    target_link_libraries(spc_lib_obj PRIVATE boost)
else ()
    target_compile_definitions(spc_lib_obj PUBLIC NUMCPP_NO_USE_BOOST)
endif ()

target_compile_features(spc_lib_obj PUBLIC cxx_std_14)
target_link_libraries(spc_lib_obj PRIVATE NumCpp::NumCpp)


add_library(spc_lib SHARED $<TARGET_OBJECTS:spc_lib_obj>)
add_library(spc_libStatic STATIC $<TARGET_OBJECTS:spc_lib_obj>)
set_target_properties(spc_libStatic PROPERTIES OUTPUT_NAME spc_lib)

# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
enable_testing()
add_subdirectory(tests)

