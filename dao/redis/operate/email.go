package operate

import "context"

/*
redis 中邮件地址集合的 CRUD 操作。因为邮件的地址信息需要频繁访问和更新的数据，使用 Redis 可以提高性能和响应速度。
*/

const EmailKey = "EmailKey" // email set(无序集合) 的键值

// AddEmails 向 redis set 中添加 Emails
func (r *RDB) AddEmails(ctx context.Context, emails ...string) error {
	if len(emails) == 0 {
		return nil
	}
	data := make([]interface{}, len(emails))
	for i, email := range emails {
		data[i] = email
	}
	return r.rdb.SAdd(ctx, EmailKey, data...).Err() // 向键值为 EmailKey 的集合中添加邮箱地址集合
}

// ExistEmail 检查指定的 email 是否存在于 set 中
func (r *RDB) ExistEmail(ctx context.Context, email string) (bool, error) {
	// 传递单个的 string 类型的 email 数据时，go 会自动将其转换为 interface 类型传入
	return r.rdb.SIsMember(ctx, EmailKey, email).Result()
}

// DeleteEmail 从 set 中删除指定的 email
func (r *RDB) DeleteEmail(ctx context.Context, email string) error {
	return r.rdb.SRem(ctx, EmailKey, email).Err()
}

// UpdateEmail 在 set 更新中指定的 email
func (r *RDB) UpdateEmail(ctx context.Context, oleEmail, newEmail string) error {
	if err := r.DeleteEmail(ctx, oleEmail); err != nil {
		return err
	}
	return r.rdb.SAdd(ctx, EmailKey, newEmail).Err()
}

// ReloadEmails 从 set 中重新加载 email 集合(删除 set 集合中的所有 emails，并重新添加)
func (r *RDB) ReloadEmails(ctx context.Context, emails ...string) error {
	if err := r.rdb.Del(ctx, EmailKey).Err(); err != nil {
		return err
	}
	return r.AddEmails(ctx, emails...)
}
