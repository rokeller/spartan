package server

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	DefaultPermissionsPolicy = PermissionsPolicy{
		Accelerometer:              AllowNonePermission{},
		AmbientLightSensor:         AllowNonePermission{},
		AriaNotify:                 AllowNonePermission{},
		AttributionReporting:       AllowNonePermission{},
		Autoplay:                   AllowNonePermission{},
		Bluetooth:                  AllowNonePermission{},
		BrowsingTopics:             AllowNonePermission{},
		Camera:                     AllowNonePermission{},
		CapturedSurfaceControl:     AllowNonePermission{},
		ComputePressure:            AllowNonePermission{},
		CrossOriginIsolated:        AllowNonePermission{},
		DeferredFetch:              AllowNonePermission{},
		DeferredFetchMinimal:       AllowNonePermission{},
		DisplayCapture:             AllowNonePermission{},
		EncryptedMedia:             AllowNonePermission{},
		Fullscreen:                 AllowNonePermission{},
		Gamepad:                    AllowNonePermission{},
		Geolocation:                AllowNonePermission{},
		Gyroscope:                  AllowNonePermission{},
		Hid:                        AllowNonePermission{},
		IdentityCredentialGet:      AllowNonePermission{},
		IdleDetection:              AllowNonePermission{},
		LanguageDetector:           AllowNonePermission{},
		LocalFonts:                 AllowNonePermission{},
		Magnetometer:               AllowNonePermission{},
		Microphone:                 AllowNonePermission{},
		Midi:                       AllowNonePermission{},
		OnDeviceSpeechRecognition:  AllowNonePermission{},
		OtpCredentials:             AllowNonePermission{},
		Payment:                    AllowNonePermission{},
		PictureInPicture:           AllowNonePermission{},
		PublickeyCredentialsCreate: AllowNonePermission{},
		PublickeyCredentialsGet:    AllowNonePermission{},
		ScreenWakeLock:             AllowNonePermission{},
		Serial:                     AllowNonePermission{},
		SpeakerSelection:           AllowNonePermission{},
		StorageAccess:              AllowNonePermission{},
		Translator:                 AllowNonePermission{},
		Summarizer:                 AllowNonePermission{},
		Usb:                        AllowNonePermission{},
		WebShare:                   AllowNonePermission{},
		WindowManagement:           AllowNonePermission{},
		XrSpatialTracking:          AllowNonePermission{},
	}
)

type PermissionAllowValue interface {
	Value() string
	PermissionAllowValueMarker()
}

type PermissionAllowListItem interface {
	Value() string
	PermissionAllowListItemMarker()
}

type PermissionsPolicy struct {
	Accelerometer              PermissionAllowValue
	AmbientLightSensor         PermissionAllowValue
	AriaNotify                 PermissionAllowValue
	AttributionReporting       PermissionAllowValue
	Autoplay                   PermissionAllowValue
	Bluetooth                  PermissionAllowValue
	BrowsingTopics             PermissionAllowValue
	Camera                     PermissionAllowValue
	CapturedSurfaceControl     PermissionAllowValue
	ComputePressure            PermissionAllowValue
	CrossOriginIsolated        PermissionAllowValue
	DeferredFetch              PermissionAllowValue
	DeferredFetchMinimal       PermissionAllowValue
	DisplayCapture             PermissionAllowValue
	EncryptedMedia             PermissionAllowValue
	Fullscreen                 PermissionAllowValue
	Gamepad                    PermissionAllowValue
	Geolocation                PermissionAllowValue
	Gyroscope                  PermissionAllowValue
	Hid                        PermissionAllowValue
	IdentityCredentialGet      PermissionAllowValue
	IdleDetection              PermissionAllowValue
	LanguageDetector           PermissionAllowValue
	LocalFonts                 PermissionAllowValue
	Magnetometer               PermissionAllowValue
	Microphone                 PermissionAllowValue
	Midi                       PermissionAllowValue
	OnDeviceSpeechRecognition  PermissionAllowValue
	OtpCredentials             PermissionAllowValue
	Payment                    PermissionAllowValue
	PictureInPicture           PermissionAllowValue
	PublickeyCredentialsCreate PermissionAllowValue
	PublickeyCredentialsGet    PermissionAllowValue
	ScreenWakeLock             PermissionAllowValue
	Serial                     PermissionAllowValue
	SpeakerSelection           PermissionAllowValue
	StorageAccess              PermissionAllowValue
	Translator                 PermissionAllowValue
	Summarizer                 PermissionAllowValue
	Usb                        PermissionAllowValue
	WebShare                   PermissionAllowValue
	WindowManagement           PermissionAllowValue
	XrSpatialTracking          PermissionAllowValue
}

func (p *PermissionsPolicy) AddToResponse(w http.ResponseWriter) {
	w.Header().Add("Permissions-Policy", p.HeaderValue())
}

func (p *PermissionsPolicy) HeaderValue() string {
	directives := []string{}

	directives = appendPermissionDirectives(directives, p.Accelerometer, "accelerometer")
	directives = appendPermissionDirectives(directives, p.AmbientLightSensor, "ambient-light-sensor")
	directives = appendPermissionDirectives(directives, p.AriaNotify, "aria-notify")
	directives = appendPermissionDirectives(directives, p.AttributionReporting, "attribution-reporting")
	directives = appendPermissionDirectives(directives, p.Autoplay, "autoplay")
	directives = appendPermissionDirectives(directives, p.Bluetooth, "bluetooth")
	directives = appendPermissionDirectives(directives, p.BrowsingTopics, "browsing-topics")
	directives = appendPermissionDirectives(directives, p.Camera, "camera")
	directives = appendPermissionDirectives(directives, p.CapturedSurfaceControl, "captured-surface-control")
	directives = appendPermissionDirectives(directives, p.ComputePressure, "compute-pressure")
	directives = appendPermissionDirectives(directives, p.CrossOriginIsolated, "cross-origin-isolated")
	directives = appendPermissionDirectives(directives, p.DeferredFetch, "deferred-fetch")
	directives = appendPermissionDirectives(directives, p.DeferredFetchMinimal, "deferred-fetch-minimal")
	directives = appendPermissionDirectives(directives, p.DisplayCapture, "display-capture")
	directives = appendPermissionDirectives(directives, p.EncryptedMedia, "encrypted-media")
	directives = appendPermissionDirectives(directives, p.Fullscreen, "fullscreen")
	directives = appendPermissionDirectives(directives, p.Gamepad, "gamepad")
	directives = appendPermissionDirectives(directives, p.Geolocation, "geolocation")
	directives = appendPermissionDirectives(directives, p.Gyroscope, "gyroscope")
	directives = appendPermissionDirectives(directives, p.Hid, "hid")
	directives = appendPermissionDirectives(directives, p.IdentityCredentialGet, "identity-credential-get")
	directives = appendPermissionDirectives(directives, p.IdleDetection, "idle-detection")
	directives = appendPermissionDirectives(directives, p.LanguageDetector, "language-detector")
	directives = appendPermissionDirectives(directives, p.LocalFonts, "local-fonts")
	directives = appendPermissionDirectives(directives, p.Magnetometer, "magnetometer")
	directives = appendPermissionDirectives(directives, p.Microphone, "microphone")
	directives = appendPermissionDirectives(directives, p.Midi, "midi")
	directives = appendPermissionDirectives(directives, p.OnDeviceSpeechRecognition, "on-device-speech-recognition")
	directives = appendPermissionDirectives(directives, p.OtpCredentials, "otp-credentials")
	directives = appendPermissionDirectives(directives, p.Payment, "payment")
	directives = appendPermissionDirectives(directives, p.PictureInPicture, "picture-in-picture")
	directives = appendPermissionDirectives(directives, p.PublickeyCredentialsCreate, "publickey-credentials-create")
	directives = appendPermissionDirectives(directives, p.PublickeyCredentialsGet, "publickey-credentials-get")
	directives = appendPermissionDirectives(directives, p.ScreenWakeLock, "screen-wake-lock")
	directives = appendPermissionDirectives(directives, p.Serial, "serial")
	directives = appendPermissionDirectives(directives, p.SpeakerSelection, "speaker-selection")
	directives = appendPermissionDirectives(directives, p.StorageAccess, "storage-access")
	directives = appendPermissionDirectives(directives, p.Translator, "translator")
	directives = appendPermissionDirectives(directives, p.Summarizer, "summarizer")
	directives = appendPermissionDirectives(directives, p.Usb, "usb")
	directives = appendPermissionDirectives(directives, p.WebShare, "web-share")
	directives = appendPermissionDirectives(directives, p.WindowManagement, "window-management")
	directives = appendPermissionDirectives(directives, p.XrSpatialTracking, "xr-spatial-tracking")

	if len(directives) > 0 {
		return strings.Join(directives, "; ")
	}
	return ""
}

func appendPermissionDirectives(directives []string, a PermissionAllowValue, name string) []string {
	if nil != a {
		return appendPermissionDirectiveValue(directives, name, a.Value())
	}

	return directives
}
func appendPermissionDirectiveValue(directives []string, name, value string) []string {
	var directive string
	if value != "" {
		directive = fmt.Sprintf("%s=%s", name, value)
	} else {
		directive = name
	}
	return append(directives, directive)
}

type AllowWildcardPermission struct{}

var _ PermissionAllowValue = AllowWildcardPermission{}

func (AllowWildcardPermission) Value() string {
	return "*"
}
func (AllowWildcardPermission) PermissionAllowValueMarker() {}

type AllowNonePermission struct{}

var _ PermissionAllowValue = AllowNonePermission{}

func (AllowNonePermission) Value() string {
	return "()"
}
func (AllowNonePermission) PermissionAllowValueMarker() {}

type AllowMultiplePermission []PermissionAllowListItem

var _ PermissionAllowValue = AllowMultiplePermission{}

func (a AllowMultiplePermission) Value() string {
	values := make([]string, len(a))
	for i, allowed := range a {
		values[i] = allowed.Value()
	}

	return fmt.Sprintf("(%s)", strings.Join(values, " "))
}
func (AllowMultiplePermission) PermissionAllowValueMarker() {}

type AllowSelfPermission struct{}

var _ PermissionAllowListItem = AllowSelfPermission{}

func (AllowSelfPermission) Value() string {
	return "self"
}
func (AllowSelfPermission) PermissionAllowListItemMarker() {}

type AllowSrcPermission struct{}

var _ PermissionAllowListItem = AllowSrcPermission{}

func (AllowSrcPermission) Value() string {
	return "src"
}
func (AllowSrcPermission) PermissionAllowListItemMarker() {}

type AllowOriginPermission string

var _ PermissionAllowListItem = AllowOriginPermission("")

func (a AllowOriginPermission) Value() string {
	return fmt.Sprintf("%q", a)
}
func (AllowOriginPermission) PermissionAllowListItemMarker() {}
