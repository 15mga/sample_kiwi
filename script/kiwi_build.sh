go install ../../kiwi_tool/cmd/protoc_kiwi/protoc-gen-kiwi.go
go install ../../kiwi_tool/cmd/protoc_mgo_bson/protoc-mgo-bson.go

DIR=`dirname $0`

NAME=game #项目go.mod中指定的模组名
GOOGLEPBIDR=$GOPATH/pb/ #google/protobuf的根目录，包含any.proto,api.proto等
PBDIR=$DIR/../proto/msg #项目protobuf文件根目录
KIWIDIR=../../../ #kiwi.proto存放根目录
OUTDIR=$DIR/../ #输出目录根目录
PBOUTDIR=$OUTDIR/proto #*.pb.go输出根目录
KIWIOUTDIR=$OUTDIR/internal/ #生成文件根目录

echo complie kiwi

protoc \
  --proto_path=$GOOGLEPBIDR \
  --proto_path=$KIWIDIR \
  --proto_path=$PBDIR \
  --go_out=$PBOUTDIR \
  --kiwi_out=-m=$NAME,-r=player,-db=mgo:$KIWIOUTDIR \
  $PBDIR/model/*.proto \
  $PBDIR/service/*.proto \
  $PBDIR/client/*.proto
echo kiwi finished

echo mgo bson
protoc-mgo-bson -d=$PBOUTDIR/pb
echo mgo bson finished

