package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

func TestImageFormat_String(t *testing.T) {
	testCases := []struct {
		name     string
		format   ImageFormat
		expected string
	}{
		{
			name:     "OK: jpeg",
			format:   ImageFormatJPEG,
			expected: "jpg",
		},
		{
			name:     "OK: gif",
			format:   ImageFormatGIF,
			expected: "gif",
		},
		{
			name:     "OK: png",
			format:   ImageFormatPNG,
			expected: "png",
		},
		{
			name:     "NG: initial value",
			format:   ImageFormat(0),
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.format.String()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestImageFormat_UnmarshalText(t *testing.T) {
	testCases := []struct {
		name        string
		text        []byte
		expected    ImageFormat
		expectedErr error
	}{
		{
			name:     "OK: jpg",
			text:     []byte("jpg"),
			expected: ImageFormatJPEG,
		},
		{
			name:     "OK: jpeg",
			text:     []byte("jpeg"),
			expected: ImageFormatJPEG,
		},
		{
			name:     "OK: gif",
			text:     []byte("gif"),
			expected: ImageFormatGIF,
		},
		{
			name:     "OK: png",
			text:     []byte("png"),
			expected: ImageFormatPNG,
		},
		{
			name:        "NG: empty",
			expectedErr: xerrors.New("invalid ImageFormat"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual ImageFormat
			err := actual.UnmarshalText(tc.text)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
