package utils

import (
	"time"

	"github.com/retawsolit/WeMeet-protocol/wemeet"
)

func PrepareDefaultRoomFeatures(r *wemeet.CreateRoomReq) {
	rf := r.Metadata.RoomFeatures

	if rf.RecordingFeatures == nil {
		rf.RecordingFeatures = &wemeet.RecordingFeatures{
			IsAllow:                  true,
			IsAllowCloud:             true,
			IsAllowLocal:             true,
			EnableAutoCloudRecording: false,
		}
	}

	if rf.ChatFeatures == nil {
		rf.ChatFeatures = &wemeet.ChatFeatures{
			AllowChat:       false,
			AllowFileUpload: false,
		}
	}

	if rf.SharedNotePadFeatures == nil {
		rf.SharedNotePadFeatures = &wemeet.SharedNotePadFeatures{
			AllowedSharedNotePad: false,
			IsActive:             false,
			Visible:              false,
		}
	}

	if rf.WhiteboardFeatures == nil {
		rf.WhiteboardFeatures = &wemeet.WhiteboardFeatures{
			AllowedWhiteboard: false,
			Visible:           false,
			WhiteboardFileId:  "default",
			FileName:          "default",
			TotalPages:        10,
		}
	}

	if rf.ExternalMediaPlayerFeatures == nil {
		rf.ExternalMediaPlayerFeatures = &wemeet.ExternalMediaPlayerFeatures{
			AllowedExternalMediaPlayer: false,
			IsActive:                   false,
		}
	}

	if rf.WaitingRoomFeatures == nil {
		rf.WaitingRoomFeatures = &wemeet.WaitingRoomFeatures{
			IsActive:       false,
			WaitingRoomMsg: "",
		}
	}

	if rf.BreakoutRoomFeatures == nil {
		rf.BreakoutRoomFeatures = &wemeet.BreakoutRoomFeatures{
			IsAllow:            false,
			IsActive:           false,
			AllowedNumberRooms: 6,
		}
	}

	if rf.DisplayExternalLinkFeatures == nil {
		rf.DisplayExternalLinkFeatures = &wemeet.DisplayExternalLinkFeatures{
			IsAllow:  false,
			IsActive: false,
		}
	}

	if rf.IngressFeatures == nil {
		rf.IngressFeatures = &wemeet.IngressFeatures{
			IsAllow: false,
		}
	}

	if rf.SpeechToTextTranslationFeatures == nil {
		rf.SpeechToTextTranslationFeatures = &wemeet.SpeechToTextTranslationFeatures{
			IsAllow:            false,
			IsAllowTranslation: false,
		}
	}

	if rf.EndToEndEncryptionFeatures == nil {
		rf.EndToEndEncryptionFeatures = &wemeet.EndToEndEncryptionFeatures{
			IsEnabled: false,
		}
	}

	if rf.PollsFeatures == nil {
		rf.PollsFeatures = &wemeet.PollsFeatures{
			IsAllow: rf.AllowPolls, // for backward compatibility
		}
	}

	if r.Metadata.DefaultLockSettings == nil {
		r.Metadata.DefaultLockSettings = new(wemeet.LockSettings)
	}

	r.Metadata.StartedAt = uint64(time.Now().UTC().Unix())
	r.Metadata.RoomFeatures = rf
}

func SetCreateRoomDefaultValues(r *wemeet.CreateRoomReq, maxSize, maxSizeWhiteboardFile uint64, allowedTypes []string, allowedNotepad bool) {
	rf := r.Metadata.RoomFeatures

	if rf.AutoGenUserId == nil {
		// by default, auto user id generation will be disabled
		ff := new(bool)
		rf.AutoGenUserId = ff
	}

	if rf.SharedNotePadFeatures.AllowedSharedNotePad && !allowedNotepad {
		rf.SharedNotePadFeatures.AllowedSharedNotePad = false
	}

	if rf.ChatFeatures.AllowFileUpload {
		if len(rf.ChatFeatures.AllowedFileTypes) == 0 {
			rf.ChatFeatures.AllowedFileTypes = allowedTypes
		}
		rf.ChatFeatures.MaxFileSize = &maxSize
	}

	if rf.WhiteboardFeatures.AllowedWhiteboard {
		if maxSizeWhiteboardFile < 1 {
			maxSizeWhiteboardFile = 30
		}
		rf.WhiteboardFeatures.MaxAllowedFileSize = &maxSizeWhiteboardFile
	}

	if rf.BreakoutRoomFeatures.IsAllow && rf.BreakoutRoomFeatures.AllowedNumberRooms == 0 {
		rf.BreakoutRoomFeatures.AllowedNumberRooms = 6
	}

	if rf.EndToEndEncryptionFeatures.IsEnabled {
		if !rf.EndToEndEncryptionFeatures.EnabledSelfInsertEncryptionKey {
			randomKey, err := GenerateSecureRandomStrings(32)
			if err != nil {
				randomKey = GenerateRandomStrings(32)
			}
			rf.EndToEndEncryptionFeatures.EncryptionKey = &randomKey
		}
	}
}

func SetRoomDefaultLockSettings(r *wemeet.CreateRoomReq) {
	lock := new(bool)
	if r.Metadata.DefaultLockSettings.LockScreenSharing == nil {
		*lock = true
		r.Metadata.DefaultLockSettings.LockScreenSharing = lock
	}
	if r.Metadata.DefaultLockSettings.LockWhiteboard == nil {
		*lock = true
		r.Metadata.DefaultLockSettings.LockWhiteboard = lock
	}
	if r.Metadata.DefaultLockSettings.LockSharedNotepad == nil {
		*lock = true
		r.Metadata.DefaultLockSettings.LockSharedNotepad = lock
	}
}

type RoomDefaultSettings struct {
	MaxParticipants     *uint32 `yaml:"max_participants"`
	MaxDuration         *uint64 `yaml:"max_duration"`
	MaxNumBreakoutRooms *uint32 `yaml:"max_num_breakout_rooms"`
}

func SetDefaultRoomSettings(s *RoomDefaultSettings, r *wemeet.CreateRoomReq) {
	if s == nil {
		return
	}

	if s.MaxParticipants != nil && *s.MaxParticipants > 0 {
		if r.MaxParticipants != nil {
			if *r.MaxParticipants == 0 || *r.MaxParticipants > *s.MaxParticipants {
				r.MaxParticipants = s.MaxParticipants
			}
		} else {
			r.MaxParticipants = s.MaxParticipants
		}
	}

	if s.MaxDuration != nil && *s.MaxDuration > 0 {
		if r.Metadata.RoomFeatures.RoomDuration != nil {
			if *r.Metadata.RoomFeatures.RoomDuration == 0 || *r.Metadata.RoomFeatures.RoomDuration > *s.MaxDuration {
				r.Metadata.RoomFeatures.RoomDuration = s.MaxDuration
			}
		} else {
			r.Metadata.RoomFeatures.RoomDuration = s.MaxDuration
		}
	}

	if r.EmptyTimeout == nil || *r.EmptyTimeout < 120 {
		var et uint32 = 1800 // 1800 seconds = 30 minutes
		r.EmptyTimeout = &et
	}

	// at present, we will allow max 16 rooms
	var maxNum uint32 = 16
	if s.MaxNumBreakoutRooms == nil {
		s.MaxNumBreakoutRooms = &maxNum
	} else if s.MaxNumBreakoutRooms != nil && *s.MaxParticipants > maxNum {
		s.MaxNumBreakoutRooms = &maxNum
	}

	if r.Metadata.RoomFeatures.BreakoutRoomFeatures.AllowedNumberRooms > *s.MaxNumBreakoutRooms {
		r.Metadata.RoomFeatures.BreakoutRoomFeatures.AllowedNumberRooms = *s.MaxNumBreakoutRooms
	}
}
