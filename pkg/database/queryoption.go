// Package database
// Date       : 2024/8/22 14:25
// Author     : Amu
// Description:
package database

type QueryOption func(db *DB) *DB

func OptionDB(db *DB, options ...QueryOption) *DB {
	for _, option := range options {
		db = option(db)
	}
	return db
}

func WithTable(tableName string) QueryOption {
	return func(db *DB) *DB {
		db.DB = db.DB.Table(tableName)
		return db
	}
}

func WithId(id string) QueryOption {
	return func(db *DB) *DB {
		if id == "" {
			return db
		}
		db.DB = db.DB.Where("id = ?", id)
		return db
	}
}

func WithIds(ids []string) QueryOption {
	return func(db *DB) *DB {
		if len(ids) == 0 {
			return db
		}
		db.DB = db.DB.Where("id IN (?)", ids)
		return db
	}
}

func WithUserName(username string) QueryOption {
	return func(db *DB) *DB {
		if username == "" {
			return db
		}
		db.DB = db.DB.Where("username = ?", username)
		return db
	}
}

func WithUsernames(usernames []string) QueryOption {
	return func(db *DB) *DB {
		if len(usernames) == 0 {
			return db
		}
		db.DB = db.DB.Where("username IN (?)", usernames)
		return db
	}
}

func WithName(name string) QueryOption {
	return func(db *DB) *DB {
		if name == "" {
			return db
		}
		db.DB = db.DB.Where("name = ?", name)
		return db
	}
}

func WithNames(names []string) QueryOption {
	return func(db *DB) *DB {
		if len(names) == 0 {
			return db
		}
		db.DB = db.DB.Where("name IN (?)", names)
		return db
	}
}

func WithStatus(status int) QueryOption {
	return func(db *DB) *DB {
		if status == 0 {
			return db
		}
		db.DB = db.DB.Where("status = ?", status)
		return db
	}
}

func WithOffset(offset int) QueryOption {
	if offset < 0 {
		return nil
	}
	return func(db *DB) *DB {
		db.DB = db.DB.Offset(offset)
		return db
	}
}

func WithLimit(limit int) QueryOption {
	if limit <= 0 {
		return nil
	}
	return func(db *DB) *DB {
		db.DB = db.DB.Limit(limit)
		return db
	}
}

func OrderBy(value interface{}) QueryOption {
	return func(db *DB) *DB {
		db.DB = db.DB.Order(value)
		return db
	}
}
