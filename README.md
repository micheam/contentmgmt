# imgcontent

画像コンテンツ管理用 CLI ツールです。  
ブログ記事などに使用する画像ファイルをクラウド環境のストレージに登録・参照する機能を提供します。

初期開発時点では、[Google Cloud Storage]をターゲットとしています。  

## Installation

[Notice] We have not yet distributed prebuilt binaries 😴

clone this repository:

    git clone https://github.com/micheam/contentmgmt

go to cmd directory:

    cd ./contentmgmt/cmd/imgcontent && go install

`imgcontent` binary will be installed:

    imgcontent help

## Environment values

env | content | note
--|--|--
IMGCONTENT_GCS_CREDENTIALS | Path to Google CloudStorage Credential file (json) | $HOME/.config/imgcontent/your-bucket-999999999999.json
IMGCONTENT_GCS_BUCKET | bucket name | your-bucket

## Usage

```
NAME:
   imagecontent - manage img content

USAGE:
   imgcontent [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Michto Maeda <https://github.com/micheam>

COMMANDS:
     upload   upload file as a web content
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Command: Upload

```
NAME:
   imgcontent upload - upload file as a web content

USAGE:
   imgcontent upload [command options] <filepath>

OPTIONS:
   --format value  Display result with specified format. [mkd,html,adoc]
   --clipboard, -c  Write result to clipboard
```

Upload your content to Google-CloudStorage:

    imgcontent upload ./sample.jpeg

Result will print via stdout:

```console
$ imgcontent upload ./sample.jpeg
https://storage.googleapis.com/micheam-image-content/2019/11/17/070934.image.jpg
```

You can specify result format with `--format`:

```console
$ imgcontent upload --format mkd ./sample.jpeg
![sample.jpeg](https://storage.googleapis.com/micheam-image-content/2019/11/19/100607.sample.jpeg)
```

## Author
[micheam](https://github.com/micheam) - <michito.maeda@gmail.com>

[Google Cloud Storage]: https://cloud.google.com/storage/
