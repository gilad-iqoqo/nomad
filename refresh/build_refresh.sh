set -x 
set -e

BASEDIR=$(dirname "$0")
echo "$BASEDIR"


$BASEDIR/flush.sh
$BASEDIR/build.sh
$BASEDIR/refresh.sh
