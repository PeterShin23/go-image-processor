package domain

import "testing"

// TestMediaProfileValidate is a table-driven test skeleton. Fill in cases and
// remove the t.Skip once you implement MediaProfile.Validate.
//
// Aim to cover: a valid profile, zero/negative width, zero/negative height,
// quality below 1, quality above 100, an unknown format, and an unknown resize
// mode.
func TestMediaProfileValidate(t *testing.T) {
	t.Skip("TODO (story 02): implement MediaProfile.Validate and fill these cases")

	tests := []struct {
		name    string
		profile MediaProfile
		wantErr bool
	}{
		// {name: "valid jpeg", profile: MediaProfile{Name: "thumb", Width: 200, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}, wantErr: false},
		// TODO: add the invalid cases described above.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.profile.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateProfilesRejectsDuplicates should assert that two profiles sharing
// a Name are rejected. Implement ValidateProfiles first.
func TestValidateProfilesRejectsDuplicates(t *testing.T) {
	t.Skip("TODO (story 02): implement ValidateProfiles and this test")
}
