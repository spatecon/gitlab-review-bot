package en_EN

import (
	"testing"

	"github.com/gertd/go-pluralize"
	"github.com/stretchr/testify/require"
)

func TestTools_Plural(t *testing.T) {
	t.Parallel()

	client := pluralize.NewClient()

	type args struct {
		n         int
		wordForms []string
	}

	tests := []struct {
		name       string
		pluralizer *pluralize.Client
		args       args
		want       string
	}{
		{
			name:       "singular",
			pluralizer: client,
			args: args{
				n:         1,
				wordForms: []string{"form"},
			},
			want: "form",
		},
		{
			name:       "plural",
			pluralizer: client,
			args: args{
				n:         2,
				wordForms: []string{"form"},
			},
			want: "forms",
		},
		{
			name:       "zero",
			pluralizer: client,
			args: args{
				n:         0,
				wordForms: []string{"form"},
			},
			want: "forms",
		},
		{
			name:       "difficult word",
			pluralizer: client,
			args: args{
				n:         0,
				wordForms: []string{"person"},
			},
			want: "people",
		},
		{
			name:       "another difficult word",
			pluralizer: client,
			args: args{
				n:         0,
				wordForms: []string{"child"},
			},
			want: "children",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tools := &Tools{
				pluralizer: tt.pluralizer,
			}

			actual := tools.Plural(tt.args.n, tt.args.wordForms...)

			require.Equal(t, tt.want, actual)
		})
	}
}
