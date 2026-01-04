package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string,maxOpenConn int,maxIdlConn int,maxIdlTime string)( *sql.DB,error){
	db,err:=sql.Open("postgres",addr)

	if err!=nil {
		return nil,err
	}
	db.SetMaxOpenConns(maxOpenConn)

	duration,err:=time.ParseDuration(maxIdlTime)
	if err!=nil{
		return nil,err
	}
	db.SetConnMaxIdleTime(duration)
	db.SetMaxIdleConns(maxIdlConn)

	ctx,cancel:= context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err!=nil{
		return nil ,err
	}
	return db,nil
}