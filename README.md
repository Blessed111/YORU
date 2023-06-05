# YORU

run database 

docker run -it -p 5432:5432 -e POSTGRES_USER=postgres  -e POSTGRES_DB=yoru -e POSTGRES_PASSWORD=1111 --name database postgres