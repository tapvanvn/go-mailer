DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

arg_num=$#
if [ $arg_num -ne 1 ];then 
    echo "Invalid call"
    echo "syntax : ./build.sh <alpine | rpi64>"
    exit 1
fi
arch=$1
pushd "$DIR/../"

tag=$(<./version.txt)

server_url=tapvanvn

docker build -t $server_url/${arch}_mailer:$tag  -f docker/$arch.dockerfile ./

docker push $server_url/${arch}_mailer:$tag

popd