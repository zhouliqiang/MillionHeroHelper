## 百万英雄答题助手使用说明

Linux用户下载/binary/linux/main使用，macOS用户下载/binary/mac/main使用。只支持屏幕是1920*1080分辨率的Android手机或相同分辨率的Genymotion。自己使用时粗略统计Linux上执行一次约在1.5s～2.9s左右，macOS上略慢些，约在2s～3.5s左右。

### 例子

题目：

![quiz](http://omohqogal.bkt.clouddn.com/quiz.png)

程序在Linux上的运行情况：

![mac_result](http://omohqogal.bkt.clouddn.com/linux_result.png)

程序在macOS上的运行情况：

![mac_result](http://omohqogal.bkt.clouddn.com/mac_result.jpg)


### 使用说明

在运行程序之前要确保：
1. 保持手机与电脑USB线连接。
2. 手机打开USB调试（设置 -> 开发者选项 -> USB调试），确保adb命令正常。

为了使程序运行起来，需要在你的系统中安装用到的依赖库。分别介绍Linux和macOS。

#### Linux

因为代码大部分是在Fedora Linux上写的，先介绍给使用Linux的朋友们。

###### 1.安装Golang

小助手是用Golang写的，运行它需要一个Golang环境。访问Golang官方网站（golang.org）下载最新版Golang Linux版本。

> https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz

如果懒得看文档，可以执行下面的命令进行安装

    tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz

然后把`/usr/local/go/bin`加入到PATH环境变量，通常情况下是在`/etc/profile`或者`$HOME/.profile`或者`$HOME/.bashrc`或者`$HOME/.zshrc`中添加

    export PATH=$PATH:/usr/local/go/bin

###### 2.安装tesseract-ocr

小助手的文字识别依赖于tesseract-ocr，目前使用3.05.x版本。

Fedora Linux

    sudo dnf install tesseract
    sudo dnf install tesseract-ocr
    sudo dnf install tesseract-devel

Ubuntu Linux

    sudo apt-get install tesseract-orc
    sudo apt-get install tesseract-orc-dev
    sudo apt-get install libleptonica-dev

为了正常识别简体中文，需要训练数据，下载`chi_sim.traineddata`文件，放到`/usr/share/tesseract-ocr/tessdata`也可能是`/usr/share/tessdata`目录。

> https://github.com/tesseract-ocr/tessdata/blob/master/chi_sim.traineddata

###### 3.还需要把Android ADB Tools添加到$PATH

相信绝大部分的Android开发者都已经具备了，当然不排除有朋友忘记配置环境变量了。

    export ANDROID_HOME=$HOME/Android/Sdk
    export PATH=$PATH:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools

可以在命令行中执行

    echo $PATH
    
检查一下Golang和Android ADB Tools是否都添加PATH中了。

当题目出现时及时在终端中执行
    
    ./main

#### macOS

###### 1.安装Golang

macOS安装Golang和添加环境变量同Linux。只需下载macOS的包。

> https://dl.google.com/go/go1.9.2.darwin-amd64.pkg

###### 2.安装tesseract-ocr

macOS安装tesseract-ocr方式有两种：

###### （1）Homebrew（推荐）

    brew install tesseract

如果选择使用Homebrew安装，需下载`chi_sim.traineddata`文件放到`/usr/local/Cellar/tesseract/$VERSION/share/tessdata`目录中。

> https://github.com/tesseract-ocr/tessdata/blob/master/chi_sim.traineddata

###### （2）MacPorts

    sudo port install tesseract

如果选择使用MacPorts安装，需再执行

    sudo port install tesseract-<chi_sim>

下载简体中文训练数据。

###### 3.还需要把Android ADB Tools添加到$PATH

操作同Linux。

当题目出现时及时在终端中执行
    
    ./main

#### Windows

由于tesseract-ocr依赖问题还没解决暂不支持。

