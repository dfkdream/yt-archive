package mpd

import (
	"reflect"
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	type args struct {
		m1 MPD
		m2 MPD
	}
	tests := []struct {
		name string
		args args
		want MPD
	}{
		{
			name: "MediaPresentationDuration",
			args: args{
				m1: MPD{
					MediaPresentationDuration: ISO8601Duration(1 * time.Second),
				},
				m2: MPD{
					MediaPresentationDuration: ISO8601Duration(2 * time.Second),
				},
			},
			want: MPD{
				MediaPresentationDuration: ISO8601Duration(2 * time.Second),
			},
		},
		{
			name: "Lang",
			args: args{
				m1: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType:       "video/webm",
								Representation: []Representation{},
							},
						},
					},
				},
				m2: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType:       "video/webm",
								Lang:           "eng",
								Representation: []Representation{},
							},
						},
					},
				},
			},
			want: MPD{
				Period: Period{
					AdaptationSet: []AdaptationSet{
						{
							MimeType:       "video/webm",
							Lang:           "eng",
							Representation: []Representation{},
						},
					},
				},
			},
		},
		{
			name: "audio + audio",
			args: args{
				m1: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "audio/webm",
								Representation: []Representation{
									{
										BaseURL: "test 1",
									},
								},
							},
						},
					},
				},
				m2: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "audio/webm",
								Representation: []Representation{
									{
										BaseURL: "test 2",
									},
								},
							},
						},
					},
				},
			},
			want: MPD{
				Period: Period{
					AdaptationSet: []AdaptationSet{
						{
							MimeType: "audio/webm",
							Representation: []Representation{
								{
									Id:      0,
									BaseURL: "test 1",
								},
								{
									Id:      1,
									BaseURL: "test 2",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "audio + video",
			args: args{
				m1: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "audio/webm",
								Representation: []Representation{
									{
										BaseURL: "test 1",
									},
								},
							},
						},
					},
				},
				m2: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "video/webm",
								Width:    10,
								Height:   10,
								Representation: []Representation{
									{
										BaseURL: "test 2",
									},
								},
							},
						},
					},
				},
			},
			want: MPD{
				Period: Period{
					AdaptationSet: []AdaptationSet{
						{
							MimeType: "video/webm",
							Id:       0,
							Representation: []Representation{
								{
									Id:      0,
									Width:   10,
									Height:  10,
									BaseURL: "test 2",
								},
							},
						},
						{
							MimeType: "audio/webm",
							Id:       1,
							Representation: []Representation{
								{
									Id:      1,
									BaseURL: "test 1",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "video + audio",
			args: args{
				m1: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "video/webm",
								Width:    10,
								Height:   10,
								Representation: []Representation{
									{
										BaseURL: "test 1",
									},
								},
							},
						},
					},
				},
				m2: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "audio/webm",
								Representation: []Representation{
									{
										BaseURL: "test 2",
									},
								},
							},
						},
					},
				},
			},
			want: MPD{
				Period: Period{
					AdaptationSet: []AdaptationSet{
						{
							MimeType: "video/webm",
							Id:       0,
							Representation: []Representation{
								{
									Id:      0,
									Width:   10,
									Height:  10,
									BaseURL: "test 1",
								},
							},
						},
						{
							MimeType: "audio/webm",
							Id:       1,
							Representation: []Representation{
								{
									Id:      1,
									BaseURL: "test 2",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "video + video",
			args: args{
				m1: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "video/webm",
								Width:    10,
								Height:   10,
								Representation: []Representation{
									{
										BaseURL: "test 1",
									},
								},
							},
						},
					},
				},
				m2: MPD{
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "video/webm",
								Width:    10,
								Height:   10,
								Representation: []Representation{
									{
										BaseURL: "test 2",
									},
								},
							},
						},
					},
				},
			},
			want: MPD{
				Period: Period{
					AdaptationSet: []AdaptationSet{
						{
							MimeType: "video/webm",
							Id:       0,
							Representation: []Representation{
								{
									Id:      0,
									Width:   10,
									Height:  10,
									BaseURL: "test 1",
								},
								{
									Id:      1,
									Width:   10,
									Height:  10,
									BaseURL: "test 2",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "mix",
			args: args{
				m1: MPD{
					MediaPresentationDuration: ISO8601Duration(100 * time.Second),
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "video/webm",
								Width:    10,
								Height:   10,
								Representation: []Representation{
									{
										BaseURL: "test video 1",
									},
								},
							},
							{
								MimeType: "audio/webm",
								Representation: []Representation{
									{
										BaseURL: "test audio 1",
									},
								},
							},
						},
					},
				},
				m2: MPD{
					MediaPresentationDuration: ISO8601Duration(100.1 * float64(time.Second)),
					Period: Period{
						AdaptationSet: []AdaptationSet{
							{
								MimeType: "video/webm",
								Representation: []Representation{
									{
										Width:   20,
										Height:  20,
										BaseURL: "test video 2",
									},
									{
										Width:   30,
										Height:  30,
										BaseURL: "test video 3",
									},
								},
							},
						},
					},
				},
			},
			want: MPD{
				MediaPresentationDuration: ISO8601Duration(100.1 * float64(time.Second)),
				Period: Period{
					AdaptationSet: []AdaptationSet{
						{
							MimeType: "video/webm",
							Id:       0,
							Representation: []Representation{
								{
									Id:      0,
									Width:   10,
									Height:  10,
									BaseURL: "test video 1",
								},
								{
									Id:      1,
									Width:   20,
									Height:  20,
									BaseURL: "test video 2",
								},
								{
									Id:      2,
									Width:   30,
									Height:  30,
									BaseURL: "test video 3",
								},
							},
						},
						{
							MimeType: "audio/webm",
							Id:       1,
							Representation: []Representation{
								{
									Id:      3,
									BaseURL: "test audio 1",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.m1, tt.args.m2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
