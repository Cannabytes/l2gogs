package models

import (
	"context"
	"l2gogameserver/data/logger"
	"l2gogameserver/db"
)

func (c *Client) SaveUser() {
	c.saveLocation()
	c.saveOnlineTime()
}

func (c *Client) saveLocation() {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()

	sql := `UPDATE "characters" SET "x" = $1, "y" = $2, "z" = $3 WHERE "object_id" = $4`
	x, y, z := c.CurrentChar.GetXYZ()
	_, err = dbConn.Exec(context.Background(), sql, x, y, z, c.CurrentChar.ObjectID())
	if err != nil {
		logger.Error.Panicln(err)
	}
}

func (c *Client) saveOnlineTime() {
	dbConn, err := db.GetConn()
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer dbConn.Release()
	sql := `UPDATE "characters" SET "online_time" = $1 WHERE "object_id" = $2`
	_, err = dbConn.Exec(context.Background(), sql, c.CurrentChar.OnlineTime(), c.CurrentChar.ObjectID())
	if err != nil {
		logger.Error.Panicln(err)
	}

}
