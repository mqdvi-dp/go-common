package constants

type StatusVerified string

const (
	Full       StatusVerified = "FULL_VERIFIED"
	Verified   StatusVerified = "VERIFIED"
	Manual     StatusVerified = "MANUAL_VERIFICATION"
	Accepted   StatusVerified = "ACCEPTED"
	Trial      StatusVerified = "TRIAL_MODE"
	Unverified StatusVerified = "UNVERIFIED"
)

func (sv StatusVerified) IsValid() bool {
	return sv == Full || sv == Verified || sv == Manual || sv == Accepted || sv == Trial
}

func (sv StatusVerified) IsTrial() bool {
	if !sv.IsValid() {
		return true
	}

	if sv == "" {
		return true
	}

	return sv == Trial
}

func (sv StatusVerified) IsUnverified() bool {
	if !sv.IsValid() {
		return true
	}

	if sv == "" {
		return true
	}

	return sv == Unverified
}

func (sv StatusVerified) String() string {
	return string(sv)
}
