package domain

import "testing"

// TestMediaProfileValidate is a table-driven test skeleton. Fill in cases and
// remove the t.Skip once you implement MediaProfile.Validate.
//
// Aim to cover: a valid profile, zero/negative width, zero/negative height,
// quality below 1, quality above 100, an unknown format, and an unknown resize
// mode.
func TestMediaProfileValidate(t *testing.T) {
	// t.Skip("TODO (story 02): implement MediaProfile.Validate and fill these cases")

	tests := []struct {
		name    string
		profile MediaProfile
		wantErr bool
	}{
		{name: "valid jpeg crop", profile: MediaProfile{Name: "thumb", Width: 200, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}},
		{name: "valid png fit", profile: MediaProfile{Name: "large", Width: 800, Height: 600, ResizeMode: ResizeFit, Format: FormatPNG, Quality: 100}},
		{name: "zero width", profile: MediaProfile{Name: "thumb", Width: 0, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}, wantErr: true},
		{name: "negative width", profile: MediaProfile{Name: "thumb", Width: -1, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}, wantErr: true},
		{name: "zero height", profile: MediaProfile{Name: "thumb", Width: 200, Height: 0, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}, wantErr: true},
		{name: "negative height", profile: MediaProfile{Name: "thumb", Width: 200, Height: -1, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}, wantErr: true},
		{name: "quality below minimum", profile: MediaProfile{Name: "thumb", Width: 200, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 0}, wantErr: true},
		{name: "quality above maximum", profile: MediaProfile{Name: "thumb", Width: 200, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 101}, wantErr: true},
		{name: "unknown resize mode", profile: MediaProfile{Name: "thumb", Width: 200, Height: 200, ResizeMode: ResizeMode("stretch"), Format: FormatJPEG, Quality: 80}, wantErr: true},
		{name: "unknown output format", profile: MediaProfile{Name: "thumb", Width: 200, Height: 200, ResizeMode: ResizeCrop, Format: OutputFormat("webp"), Quality: 80}, wantErr: true},
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
func TestValidateProfiles(t *testing.T) {
	valid := MediaProfile{Name: "thumb", Width: 200, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}

	tests := []struct {
		name     string
		profiles []MediaProfile
		wantErr  bool
	}{
		{name: "empty"},
		{name: "unique names", profiles: []MediaProfile{valid, {Name: "large", Width: 800, Height: 600, ResizeMode: ResizeFit, Format: FormatPNG, Quality: 90}}},
		{name: "duplicate names", profiles: []MediaProfile{valid, valid}, wantErr: true},
		{name: "invalid profile", profiles: []MediaProfile{{Name: "broken", Width: 0, Height: 200, ResizeMode: ResizeCrop, Format: FormatJPEG, Quality: 80}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProfiles(tt.profiles)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateProfiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
