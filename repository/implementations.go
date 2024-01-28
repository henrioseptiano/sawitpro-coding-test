package repository

import (
	"context"
	"errors"
	"log"
)

func (r *Repository) RegisterUser(input User) error {
	_, err := r.Db.Exec("INSERT INTO users ("+
		"user_id, full_name, phone_number, password, successfull_login_attempts, last_login,"+
		"created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", input.UserID, input.FullName,
		input.PhoneNumber, input.Password, input.SuccessfullLoginAttempts, input.LastLogin,
		input.CreatedAt, input.UpdatedAt)
	if err != nil {
		log.Println(err)
		return errors.New("cannot Register the user")
	}
	return err
}

// used for login
func (r *Repository) CheckUser(ctx context.Context, phoneNumber string) (*User, error) {
	output := User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, user_id, full_name, phone_number, password,"+
		" successfull_login_attempts, last_login FROM users WHERE phone_number = $1", phoneNumber).
		Scan(&output.ID, &output.UserID, &output.FullName, &output.PhoneNumber, &output.Password, &output.SuccessfullLoginAttempts,
			&output.LastLogin)
	if err != nil {
		log.Println(err)
		return nil, errors.New("there is problem in our system when performing query. please wait")
	}
	return &output, nil
}

func (r *Repository) UpdateLoginUser(ctx context.Context, input User) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE users SET successfull_login_attempts = $1, last_login = $2"+
		" WHERE id = $3", input.SuccessfullLoginAttempts, input.LastLogin, input.ID)
	if err != nil {
		log.Println(err)
		return errors.New("there is problem in our system when performing login. please wait")
	}
	return err
}

func (r *Repository) GetUserByUserId(ctx context.Context, userID string) (*User, error) {
	output := User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, user_id, full_name, phone_number, password,"+
		" successfull_login_attempts, last_login FROM users WHERE user_id = $1", userID).
		Scan(&output.ID, &output.UserID, &output.FullName, &output.PhoneNumber, &output.Password, &output.SuccessfullLoginAttempts,
			&output.LastLogin)
	if err != nil {
		log.Println(err)
		return nil, errors.New("there is problem in our system when performing query. please wait")

	}
	return &output, nil
}

func (r *Repository) UpdateUserProfile(ctx context.Context, input User) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE users SET "+
		"phone_number = $1, full_name = $2"+
		" WHERE id = $3", input.PhoneNumber, input.FullName, input.ID)
	if err != nil {
		log.Println(err)
		return errors.New("there is problem in our system when updating profile. please wait")
	}
	return err
}

func (r *Repository) CheckPhoneNumber(ctx context.Context, phoneNumber string) (int64, error) {
	count := 0
	err := r.Db.QueryRowContext(ctx, "SELECT count(id) FROM users WHERE phone_number = $1", phoneNumber).
		Scan(&count)
	if err != nil {
		log.Println(err)
		return 0, errors.New("there is problem in our system when performing query. please wait")
	}
	return int64(count), nil
}

/*func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}*/
