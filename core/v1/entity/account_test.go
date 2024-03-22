package entity

import (
	"dot/core"
	"reflect"
	"testing"
	"time"
)

func TestAccount_CheckPasword(t *testing.T) {
	account := &Account{Password: "password"}
	account.SetPassword()

	type fields struct {
		Id               string
		Otp              string
		Username         string
		Email            string
		Password         string
		PhoneNumber      string
		Age              int
		Gender           int
		Verified         int
		Metadata         map[string]any
		VerifiedDate     *time.Time
		ResetToken       string
		ResetTokenExpire *time.Time
		SentAccess       *time.Time
		Created          *time.Time
		Modified         *time.Time
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *core.CustomError
	}{
		{
			name: "password must be passed",
			fields: fields{
				Password: account.Password,
			},
			args: args{
				password: "password",
			},
			want: nil,
		},
		{
			name: "password must be invalid",
			fields: fields{
				Password: "",
			},
			args: args{},
			want: &core.CustomError{
				Code:    core.WRONG_PASSWORD,
				Message: "wrong password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Id:               tt.fields.Id,
				Otp:              tt.fields.Otp,
				Username:         tt.fields.Username,
				Email:            tt.fields.Email,
				Password:         tt.fields.Password,
				PhoneNumber:      tt.fields.PhoneNumber,
				Age:              tt.fields.Age,
				Gender:           tt.fields.Gender,
				Verified:         tt.fields.Verified,
				Metadata:         tt.fields.Metadata,
				VerifiedDate:     tt.fields.VerifiedDate,
				ResetToken:       tt.fields.ResetToken,
				ResetTokenExpire: tt.fields.ResetTokenExpire,
				SentAccess:       tt.fields.SentAccess,
				Created:          tt.fields.Created,
				Modified:         tt.fields.Modified,
			}
			if got := a.CheckPasword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckPasword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_RemoveResetToken(t *testing.T) {
	now := time.Now()
	type fields struct {
		Id               string
		Otp              string
		Username         string
		Email            string
		Password         string
		PhoneNumber      string
		Age              int
		Gender           int
		Verified         int
		Metadata         map[string]any
		VerifiedDate     *time.Time
		ResetToken       string
		ResetTokenExpire *time.Time
		SentAccess       *time.Time
		Created          *time.Time
		Modified         *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *Account
	}{
		{
			name: "reset token must be removed",
			fields: fields{
				ResetToken:       "reset token",
				ResetTokenExpire: &now,
			},
			want: &Account{
				ResetToken:       "",
				ResetTokenExpire: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Id:               tt.fields.Id,
				Otp:              tt.fields.Otp,
				Username:         tt.fields.Username,
				Email:            tt.fields.Email,
				Password:         tt.fields.Password,
				PhoneNumber:      tt.fields.PhoneNumber,
				Age:              tt.fields.Age,
				Gender:           tt.fields.Gender,
				Verified:         tt.fields.Verified,
				Metadata:         tt.fields.Metadata,
				VerifiedDate:     tt.fields.VerifiedDate,
				ResetToken:       tt.fields.ResetToken,
				ResetTokenExpire: tt.fields.ResetTokenExpire,
				SentAccess:       tt.fields.SentAccess,
				Created:          tt.fields.Created,
				Modified:         tt.fields.Modified,
			}
			if got := a.RemoveResetToken(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveResetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_SetNewPassword(t *testing.T) {
	account := &Account{Password: "password old"}
	account.SetPassword()
	type fields struct {
		Id               string
		Otp              string
		Username         string
		Email            string
		Password         string
		PhoneNumber      string
		Age              int
		Gender           int
		Verified         int
		Metadata         map[string]any
		VerifiedDate     *time.Time
		ResetToken       string
		ResetTokenExpire *time.Time
		SentAccess       *time.Time
		Created          *time.Time
		Modified         *time.Time
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *core.CustomError
	}{
		{
			name: "new password must be set",
			fields: fields{
				Password: account.Password,
			},
			args: args{
				password: "password new",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Id:               tt.fields.Id,
				Otp:              tt.fields.Otp,
				Username:         tt.fields.Username,
				Email:            tt.fields.Email,
				Password:         tt.fields.Password,
				PhoneNumber:      tt.fields.PhoneNumber,
				Age:              tt.fields.Age,
				Gender:           tt.fields.Gender,
				Verified:         tt.fields.Verified,
				Metadata:         tt.fields.Metadata,
				VerifiedDate:     tt.fields.VerifiedDate,
				ResetToken:       tt.fields.ResetToken,
				ResetTokenExpire: tt.fields.ResetTokenExpire,
				SentAccess:       tt.fields.SentAccess,
				Created:          tt.fields.Created,
				Modified:         tt.fields.Modified,
			}
			if got := a.SetNewPassword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetNewPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_SetPassword(t *testing.T) {
	type fields struct {
		Id               string
		Otp              string
		Username         string
		Email            string
		Password         string
		PhoneNumber      string
		Age              int
		Gender           int
		Verified         int
		Metadata         map[string]any
		VerifiedDate     *time.Time
		ResetToken       string
		ResetTokenExpire *time.Time
		SentAccess       *time.Time
		Created          *time.Time
		Modified         *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *core.CustomError
	}{
		{
			name: "existing password must be hashed",
			fields: fields{
				Password: "password",
			},
			want: nil,
		},
		{
			name: "ignore password hash when it's empty",
			fields: fields{
				Password: "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Id:               tt.fields.Id,
				Otp:              tt.fields.Otp,
				Username:         tt.fields.Username,
				Email:            tt.fields.Email,
				Password:         tt.fields.Password,
				PhoneNumber:      tt.fields.PhoneNumber,
				Age:              tt.fields.Age,
				Gender:           tt.fields.Gender,
				Verified:         tt.fields.Verified,
				Metadata:         tt.fields.Metadata,
				VerifiedDate:     tt.fields.VerifiedDate,
				ResetToken:       tt.fields.ResetToken,
				ResetTokenExpire: tt.fields.ResetTokenExpire,
				SentAccess:       tt.fields.SentAccess,
				Created:          tt.fields.Created,
				Modified:         tt.fields.Modified,
			}
			if got := a.SetPassword(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
