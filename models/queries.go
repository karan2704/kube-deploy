package models

var InsertAppQuery string = `INSERT INTO Apps (appid, owner)
								  VALUES ($1, $2)`

var InsertFileQuery string = `INSERT INTO Files (filename, content, appid)
							  VALUES ($1, $2, $3)`

var CheckUserExists string = `SELECT EXISTS (SELECT * FROM USER WHERE id = $1`

var SaveUserInfo string = `INSERT INTO Users (id, name, email)
						   VALUES ($1, $2, $3)`