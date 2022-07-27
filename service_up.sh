docker-compose -f deployment.yml down
while getopts "bk" arg
do
  case $arg in
    k)
      echo "Stopped."
      exit 0
      ;;
    b)
      docker-compose -f deployment.yml rm -f
      BUILD="--build"
      ;;
    ?)
      echo "Unknown args $arg... exit..."
      exit 1
      ;;
  esac
done

# Remove previous env var file
ENV_FILE=".env"
rm -f $ENV_FILE

# Set env ars
CUR_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
HOST_IP=`hostname -I | cut -d ' ' -f1`

# Encode env vars for docker compose
function encode_env {
  echo "$1=${!1}" >> $ENV_FILE
}

encode_env "CUR_DIR"
encode_env "HOST_IP"

# Start service
docker-compose -f deployment.yml up $BUILD -d
