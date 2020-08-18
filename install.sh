TAR_FILE=node.tar.gz
tar -czvf $TAR_FILE ./node/node ./data restart_node.sh

MV_TAR=/data/apps/guard
mv $TAR_FILE $MV_TAR

cd $MV_TAR
mv node node_bak
rm -rf ./data

tar -zxvf $TAR_FILE
mv ./node/node node_tmp
rm -rf ./node/
mv node_tmp node

sh restart_node.sh
