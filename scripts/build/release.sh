#!/bin/bash

# 设置变量
WORKSPACE_DIR=$(dirname $(dirname $(dirname $(realpath $0))))
BUILD_DIR="$WORKSPACE_DIR/build"
RELEASE_DIR="$WORKSPACE_DIR/release"
VERSION=$(jq -r '.version' < "$WORKSPACE_DIR/package.json")
RELEASE_FILE="website_$VERSION.tar.gz"

# 创建构建目录和发布目录
mkdir -p $BUILD_DIR
mkdir -p $RELEASE_DIR

# 运行构建任务
echo "Running build tasks..."
npm run build:css --prefix $WORKSPACE_DIR
npm run build:js --prefix $WORKSPACE_DIR
go build -o $BUILD_DIR/idiary.site $WORKSPACE_DIR/cmd/idiary/main.go

# 复制配置文件、静态文件和其他必要文件到构建目录
echo "Copying config, static files, and other necessary files..."
mkdir -p $BUILD_DIR/website

cp -r $WORKSPACE_DIR/configs $BUILD_DIR/
rsync -av --exclude 'src' $WORKSPACE_DIR/website/static/ $BUILD_DIR/website/static/
cp $WORKSPACE_DIR/go.mod $WORKSPACE_DIR/go.sum $BUILD_DIR/


# 打包构建文件和其他必要文件
echo "Packaging files..."
tar -czvf $RELEASE_DIR/$RELEASE_FILE -C $BUILD_DIR .

echo "Release package created at $RELEASE_DIR/$RELEASE_FILE"