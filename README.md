# To startup the MySQL database and run migrations, you can use the following commands:
```bash
docker run --name ecomm-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=admin -d mysql:latest  
docker exec -i ecomm-mysql mysql -uroot -padmin -e "CREATE DATABASE ecomm;"  
docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.18.3 -path=/db/migrations -database "mysql://root:admin@tcp(localhost:3306)/ecomm" up 
```