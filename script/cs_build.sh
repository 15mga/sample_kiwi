DIR=`dirname $0`

OUTDIR=$DIR/../
PBDIR=$DIR/../proto/msg
ETDIR=/Users/95eh/95eh.com/go
PBOUTDIR=$OUTDIR/proto
EGOUTDIR=$OUTDIR/internal/

echo complie cs

protoc \
  --proto_path=$GOPATH/pb/ \
  --proto_path=$ETDIR \
  --proto_path=$PBDIR \
  --csharp_out=$PBOUTDIR/cs \
  $PBDIR/model/*.proto $PBDIR/client/*.proto

echo cs finished