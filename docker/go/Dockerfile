#golang v1.8が準備されているimageを利用
FROM golang:1.18

#$GOPATHに/go/srcを追加．この後srcの下にアプリケーションフォルダを作成する為
ENV GOPATH $GOPATH:/go/src

#とりあえず更新
RUN apt-get update
    # apt-get upgrade

#revel，revel-cli，gorm，go-sql-driverのインストール
#revelにはORMがないので
# RUN go get github.com/revel/revel && \
#     go get github.com/revel/cmd/revel

#アプリケーション(myapp)をマウントするためのディレクトリを作成
RUN mkdir /go/src/app
