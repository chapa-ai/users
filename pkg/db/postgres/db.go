package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"log"
	"os"
	"path"
	"time"
	pb "users/grpc-users"
	"users/pkg/db/redis"
)

var db *pgx.Conn

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "postgres"
)

func InitDB() (*pgx.Conn, error) {
	if db == nil {
		conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
		d, err := pgx.Connect(context.Background(), conn)
		if err != nil {
			log.Fatalf("Unable to establish connection: %v", err)
		}
		db = d
	}
	return db, nil
}

func MigrateDB() error {
	conn, err := InitDB()
	if err != nil {
		return err
	}

	migrator, err := migrate.NewMigrator(context.Background(), conn, "schema_version")
	if err != nil {
		return err
	}

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dir := path.Dir(mydir)
	join := path.Join(dir, "pkg/db/postgres/migrations")

	err = migrator.LoadMigrations(join)
	if err != nil {
		return err
	}

	err = migrator.Migrate(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func InsertUser(name string, age string) (*pb.User, error) {

	var users *pb.User = &pb.User{}

	err := db.QueryRow(context.Background(), `insert into users(name, age) values($1, $2) RETURNING "id", "name", "age"`, &name, &age).Scan(&users.Id, &users.Name, &users.Age)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsers() (*pb.UserList, error) {
	rdb := redis.RedisConnect()

	var users_list *pb.UserList = &pb.UserList{}
	rows, err := db.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}

		val, _ := rdb.Get(context.Background(), user.Name).Result()
		if err != nil {
			return nil, err
		}
		if val == "" {
			_, err := rdb.SetNX(context.Background(), user.Name, user.Age, 60*time.Second).Result()
			if err != nil {
				return nil, err
			}
		}
		users_list.Users = append(users_list.Users, &user)
	}
	return users_list, nil
}

func DeleteUser(id int64) (*pb.DeletedUser, error) {
	var delUser *pb.DeletedUser = &pb.DeletedUser{}

	err := db.QueryRow(context.Background(), `delete from users where id = $1 RETURNING "id", "name", "age" `, id).Scan(&delUser.Id, &delUser.Name, &delUser.Age)
	if err != nil {
		return nil, err
	}

	return delUser, nil
}
