package ds

import "testing"

func TestEqualUser(t *testing.T) {
	t.Parallel()

	type args struct {
		u1 *BasicUser
		u2 *BasicUser
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal",
			args: args{
				u1: &BasicUser{
					GitLabID: 123,
					Name:     "name",
				},
				u2: &BasicUser{
					GitLabID: 123,
					Name:     "name",
				},
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				u1: &BasicUser{
					GitLabID: 123,
				},
				u2: &BasicUser{
					GitLabID: 123,
					Name:     "name",
				},
			},
			want: true,
		},
		{
			name: "both nil",
			args: args{
				u1: nil,
				u2: nil,
			},
			want: true,
		},
		{
			name: "one nil",
			args: args{
				u1: &BasicUser{
					GitLabID: 123,
				},
				u2: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualUser(tt.args.u1, tt.args.u2); got != tt.want {
				t.Errorf("EqualUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualUsers(t *testing.T) {
	t.Parallel()

	type args struct {
		u1 []*BasicUser
		u2 []*BasicUser
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal",
			args: args{
				u1: []*BasicUser{
					{
						GitLabID: 123,
						Name:     "name",
					},
				},
				u2: []*BasicUser{
					{
						GitLabID: 123,
						Name:     "name",
					},
				},
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				u1: []*BasicUser{
					{
						GitLabID: 123,
						Name:     "name",
					},
				},
				u2: []*BasicUser{
					{
						GitLabID: 123,
					},
				},
			},
			want: false,
		},
		{
			name: "both nil",
			args: args{
				u1: nil,
				u2: nil,
			},
			want: true,
		},
		{
			name: "one nil",
			args: args{
				u1: []*BasicUser{
					{
						GitLabID: 123,
						Name:     "name",
					},
				},
				u2: nil,
			},
			want: false,
		},
		{
			name: "different length",
			args: args{
				u1: []*BasicUser{
					{
						GitLabID: 123,
					},
				},
				u2: []*BasicUser{
					{
						GitLabID: 123,
					},
					{
						GitLabID: 431,
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualUsers(tt.args.u1, tt.args.u2); got != tt.want {
				t.Errorf("EqualUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
