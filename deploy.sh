serverName=serica-go
environment=develop

awsAccountId="404904371652"
region="ap-southeast-1"
ecrUrl="$awsAccountId.dkr.ecr.$region.amazonaws.com"
path="/home/ubuntu/$environment/$serverName"
fileName="$path/docker-compose.yml"

docker-compose -f "$fileName" -p "$serverName-${environment}" down -v

aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 404904371652.dkr.ecr.ap-southeast-1.amazonaws.com
docker pull 404904371652.dkr.ecr.ap-southeast-1.amazonaws.com/serica-go-dev:latest

docker-compose -f "$fileName" -p "$serverName-${environment}" pull
docker-compose -f "$fileName" -p "$serverName-${environment}" up -d