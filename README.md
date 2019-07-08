# imgcontent

画像コンテンツ管理用 CLI ツールです。  
ブログ記事などに使用する画像ファイルをクラウド環境のストレージに登録・参照する機能を提供します。

初期開発時点では、[Google Cloud Storage]をターゲットとしています。  

## Usage

画像ファイルのアップロード:

```sh
imgcontent upload $HOME/Pictures/awesome-image.png
```

アップロード済み画像の一覧表示:

```sh
imgcontent list --prefix 2019/07/01
```

## Requirements

TBD

## Installation

TBD

## Author
[micheam](https://github.com/micheam) - <michito.maeda@gmail.com>


[Google Cloud Storage]: https://cloud.google.com/storage/
