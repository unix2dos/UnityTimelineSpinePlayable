CURRENT_DIR=`dirname $0`
OLD_DIR=$CURRENT_DIR/dataIsland
cd $CURRENT_DIR
./command -path=../Assets/Resources/dataDef
rm -rf ../Assets/Resources/dataDef/dataConfig
mv dataIsland/dataConfig ../Assets/Resources/dataDef
rm -rf $OLD_DIR
