package unit

import (
	"strings"
	"testing"

	"github.com/harrybawsac/knmi-go/internal/parser"
)

func TestParseCSV(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		wantErr       bool
		errContains   string
	}{
		{
			name: "valid single record",
			input: `# KNMI data file header
# STN,YYYYMMDD,DDVEC,FHVEC,FG,FHX,FHXH,FHN,FHNH,FXX,FXXH,TG,TN,TNH,TX,TXH,T10N,T10NH,SQ,SP,Q,DR,RH,RHX,RHXH,PG,PX,PXH,PN,PNH,VVN,VVNH,VVX,VVXH,NG,UG,UX,UXH,UN,UNH,EV24
  260,20240101,  230,   45,   52,   72,   15,   31,    1,  100,   15,   85,   62,    6,  102,   14,   52,    6,   25,   28, 380,   10,   32,    8,   12,10250,10280,   12,10220,    6,   54,    7,   75,   15,    6,   88,   96,    7,   78,   14,    8
`,
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "multiple valid records",
			input: `# STN,YYYYMMDD,...
  260,20240101,  230,   45,   52,   72,   15,   31,    1,  100,   15,   85,   62,    6,  102,   14,   52,    6,   25,   28, 380,   10,   32,    8,   12,10250,10280,   12,10220,    6,   54,    7,   75,   15,    6,   88,   96,    7,   78,   14,    8
  260,20240102,  180,   38,   45,   65,   10,   25,    5,   90,   10,   90,   70,    3,  110,   15,   60,    3,   30,   35, 400,    5,   20,    5,    8,10260,10290,   10,10230,    5,   50,    5,   80,   12,    5,   85,   92,    5,   75,   12,   10
`,
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name: "empty values handled",
			input: `# STN,YYYYMMDD,...
  260,20240101,     ,   45,   52,     ,   15,   31,    1,     ,   15,   85,   62,    6,  102,   14,     ,    6,   25,   28, 380,     ,   32,    8,   12,10250,10280,   12,10220,    6,     ,    7,   75,   15,    6,   88,   96,    7,   78,   14,    8
`,
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name:          "empty input",
			input:         "",
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "only comments",
			input: `# This is a comment
# Another comment
# No data here
`,
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "wrong column count",
			input: `# STN,YYYYMMDD,...
  260,20240101,  230,   45
`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "column",
		},
		{
			name: "invalid date format",
			input: `# STN,YYYYMMDD,...
  260,2024-01-01,  230,   45,   52,   72,   15,   31,    1,  100,   15,   85,   62,    6,  102,   14,   52,    6,   25,   28, 380,   10,   32,    8,   12,10250,10280,   12,10220,    6,   54,    7,   75,   15,    6,   88,   96,    7,   78,   14,    8
`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := parser.ParseCSV(strings.NewReader(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(records) != tt.expectedCount {
				t.Errorf("expected %d records, got %d", tt.expectedCount, len(records))
			}
		})
	}
}

func TestParseCSVFieldValues(t *testing.T) {
	input := `# STN,YYYYMMDD,...
  260,20240115,  230,   45,   52,   72,   15,   31,    1,  100,   15,   85,   62,    6,  102,   14,   52,    6,   25,   28, 380,   10,   32,    8,   12,10250,10280,   12,10220,    6,   54,    7,   75,   15,    6,   88,   96,    7,   78,   14,    8
`

	records, err := parser.ParseCSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	r := records[0]

	// Check station ID
	if r.StationID != 260 {
		t.Errorf("expected station_id=260, got %d", r.StationID)
	}

	// Check date
	if r.Date.Year() != 2024 || r.Date.Month() != 1 || r.Date.Day() != 15 {
		t.Errorf("expected date=2024-01-15, got %v", r.Date)
	}

	// Check a few integer fields
	if r.DDVEC == nil || *r.DDVEC != 230 {
		t.Errorf("expected ddvec=230, got %v", r.DDVEC)
	}

	if r.TG == nil || *r.TG != 85 {
		t.Errorf("expected tg=85, got %v", r.TG)
	}

	if r.EV24 == nil || *r.EV24 != 8 {
		t.Errorf("expected ev24=8, got %v", r.EV24)
	}
}

func TestParseCSVEmptyFields(t *testing.T) {
	input := `# STN,YYYYMMDD,...
  260,20240115,     ,   45,     ,   72,   15,     ,    1,  100,   15,   85,     ,    6,  102,   14,   52,    6,   25,   28, 380,   10,   32,    8,   12,10250,10280,   12,10220,    6,   54,    7,   75,   15,    6,   88,   96,    7,   78,   14,    8
`

	records, err := parser.ParseCSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	r := records[0]

	// DDVEC should be nil (was empty)
	if r.DDVEC != nil {
		t.Errorf("expected ddvec=nil (empty), got %v", r.DDVEC)
	}

	// FHVEC should have value (was not empty)
	if r.FHVEC == nil || *r.FHVEC != 45 {
		t.Errorf("expected fhvec=45, got %v", r.FHVEC)
	}

	// FG should be nil (was empty)
	if r.FG != nil {
		t.Errorf("expected fg=nil (empty), got %v", r.FG)
	}
}
