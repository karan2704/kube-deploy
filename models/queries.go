package models

var InsertManifestFiles string = `INSERT INTO Apps (appid, owner)
								  VALUES ($1, $2) `