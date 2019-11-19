# imgcontent

画像コンテンツ管理用 CLI ツールです。  
ブログ記事などに使用する画像ファイルをクラウド環境のストレージに登録・参照する機能を提供します。

初期開発時点では、[Google Cloud Storage]をターゲットとしています。  

## Environment values

環境変数 | 内容 | 例
--|--|--
IMGCONTENT_GCS_CREDENTIALS | Google Cloud Storage Credential file (json) | $HOME/.config/imgcontent/your-bucket-999999999999.json
IMGCONTENT_GCS_BUCKET | 対象バケット名 | your-bucket

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

upload your jpeg to gcp:

    imgcontent upload ./sample.jpeg

result will print via stdout:

```console
$ imgcontent upload ./sample.jpeg
![image.jpg](https://storage.googleapis.com/micheam-image-content/2019/11/17/070934.image.jpg)
```


## Installation
clone this repository:

    git clone https://github.com/micheam/contentmgmt

go to cmd directory:

    cd ./contentmgmt/cmd/imgcontent && go install

`imgcontent` binary will be installed:

    imgcontent help

## Author
[micheam](https://github.com/micheam) - <michito.maeda@gmail.com>

[Google Cloud Storage]: https://cloud.google.com/storage/
