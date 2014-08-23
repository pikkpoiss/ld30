#!/bin/bash

chmod +w *.dylib

BINARY=ld30

install_name_tool -change /usr/local/opt/glfw3/lib/libglfw3.3.0.dylib @executable_path/libglfw3.dylib ${BINARY}
install_name_tool -change /usr/local/lib/libGLEW.1.10.0.dylib @executable_path/libGLEW.1.10.0.dylib ${BINARY}
install_name_tool -change /usr/lib/libGLEW.1.9.0.dylib @executable_path/libGLEW.1.9.0.dylib ${BINARY}
install_name_tool -change /usr/local/lib/libSDL-1.2.0.dylib @executable_path/libSDL-1.2.0.dylib ${BINARY}
install_name_tool -change /usr/local/lib/libSDL-1.2.0.dylib @executable_path/libSDL-1.2.0.dylib libSDL-1.2.0.dylib
install_name_tool -change /usr/local/lib/libSDL-1.2.0.dylib @executable_path/libSDL-1.2.0.dylib libSDL_mixer-1.2.0.dylib
install_name_tool -change /usr/local/lib/libSDL-1.2.0.dylib @executable_path/libSDL-1.2.0.dylib libSDL_image-1.2.0.dylib
install_name_tool -change /usr/local/lib/libSDL_image-1.2.0.dylib @executable_path/libSDL_image-1.2.0.dylib ${BINARY}
install_name_tool -change /usr/local/lib/libSDL_image-1.2.0.dylib @executable_path/libSDL_image-1.2.0.dylib libSDL_image-1.2.0.dylib
install_name_tool -change /usr/local/lib/libSDL_mixer-1.2.0.dylib @executable_path/libSDL_mixer-1.2.0.dylib ${BINARY}
install_name_tool -change /usr/local/lib/libSDL_mixer-1.2.0.dylib @executable_path/libSDL_mixer-1.2.0.dylib libSDL_mixer-1.2.0.dylib
install_name_tool -change /usr/local/Cellar/libvorbis/1.3.3/lib/libvorbis.0.dylib @executable_path/libvorbis.0.dylib libvorbisfile.dylib
install_name_tool -change /usr/local/Cellar/libvorbis/1.3.3/lib/libvorbis.0.dylib @executable_path/libvorbis.0.dylib libvorbisfile.3.dylib
install_name_tool -change /usr/local/lib/libvorbisfile.3.dylib @executable_path/libvorbisfile.3.dylib libvorbisfile.dylib
install_name_tool -change /usr/local/lib/libvorbisfile.3.dylib @executable_path/libvorbisfile.3.dylib libvorbisfile.3.dylib
install_name_tool -change /usr/local/lib/libogg.0.dylib @executable_path/libogg.0.dylib libvorbisfile.dylib
install_name_tool -change /usr/local/lib/libogg.0.dylib @executable_path/libogg.0.dylib libvorbisfile.3.dylib
install_name_tool -change /usr/local/lib/libogg.0.dylib @executable_path/libogg.0.dylib libvorbis.dylib
install_name_tool -change /usr/local/lib/libogg.0.dylib @executable_path/libogg.0.dylib libvorbis.0.dylib
install_name_tool -change /usr/local/lib/libvorbis.0.dylib @executable_path/libvorbis.0.dylib libvorbis.dylib
install_name_tool -change /usr/local/lib/libvorbis.0.dylib @executable_path/libvorbis.0.dylib libvorbis.0.dylib

chmod -w *.dylib
otool -L ${BINARY} libSDL*.dylib libvorbis*
