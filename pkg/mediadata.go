package pkg

import (
	"context"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
	"strconv"
	"strings"
	"time"
)

// AudioData is a struct that holds information about an audio stream
type AudioData struct {
	Codec    string  // The codec used to encode the audio (e.g. AAC, MP3, etc)
	Language string  // The language of the audio stream (e.g. English, Spanish, etc)
	Bitrate  float64 // The bitrate of the audio stream in bits per second
}

// SubtitleData is a struct that holds information about a subtitle stream
type SubtitleData struct {
	Language string // The language of the subtitle stream (e.g. English, Spanish, etc)
	Codec    string // The codec used to encode the subtitles (e.g. SRT, ASS, etc)
}

// MediaData is a struct that holds information about a media file
type MediaData struct {
	Size      float64        // The size of the media file in bytes
	Duration  float64        // The duration of the media file in seconds
	Codec     string         // The codec used to encode the video (e.g. H.264, H.265, etc)
	Audios    []AudioData    // An array of AudioData structs representing the audio streams in the media file
	Subtitles []SubtitleData // An array of SubtitleData structs representing the subtitle streams in the media file
}

// RetrieveMediaData retrieves media data from a file using ffprobe
func RetrieveMediaData(filePath string) (MediaData, error) {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Probe the file using ffprobe
	data, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		log.Println("Error retrieving media data:", err)
		return MediaData{}, err
	}

	// Initialize a new MediaData struct
	var mediaData MediaData

	// Extract the video stream information from the ffprobe data
	err = extractVideoStream(data, &mediaData)
	if err != nil {
		log.Println("Error extracting video stream:", err)
		return MediaData{}, err
	}

	// Extract the audio stream information from the ffprobe data
	err = extractAudioStream(data, &mediaData)
	if err != nil {
		log.Println("Error extracting audio stream:", err)
		return MediaData{}, err
	}

	// Extract the subtitle stream information from the ffprobe data
	extractSubtitleStream(data, &mediaData)

	// Return the completed media data
	return mediaData, nil
}

// extractVideoStream extracts video stream information from the ffprobe data and populates the provided MediaData struct
func extractVideoStream(data *ffprobe.ProbeData, mediaData *MediaData) error {
	// Parse the size of the media file
	size, err := strconv.ParseFloat(data.Format.Size, 64)
	if err != nil {
		log.Println("Error parsing size:", err)
		return err
	}
	mediaData.Size = size

	// Set the duration of the media file
	mediaData.Duration = data.Format.DurationSeconds

	// Extract the first video stream from the ffprobe data
	videoStream := data.FirstVideoStream()

	// Set the video codec for the media file
	mediaData.Codec = strings.ToUpper(videoStream.CodecName)

	// Return nil to indicate successful extraction
	return nil
}

// extractSubtitleStream extracts subtitle stream information from the ffprobe data and populates the provided MediaData struct
func extractSubtitleStream(data *ffprobe.ProbeData, mediaData *MediaData) {
	// Loop through all subtitle streams in the ffprobe data
	for _, subtitleStream := range data.StreamType(ffprobe.StreamSubtitle) {
		// Extract the language of the subtitle stream from the ffprobe data
		language, err := subtitleStream.TagList.GetString("language")
		if err != nil {
			// If the language tag is not available, try to extract the title tag
			language, err = subtitleStream.TagList.GetString("title")
			if err != nil {
				// If neither tag is available, set the language to "Unknown"
				language = "Unknown"
			}
		}

		// Add a new SubtitleData struct to the mediaData.subtitle array
		mediaData.Subtitles = append(mediaData.Subtitles, SubtitleData{
			Language: language,
			Codec:    strings.ToUpper(subtitleStream.CodecName),
		})
	}
}

// extractAudioStream extracts audio stream information from the ffprobe data and populates the provided MediaData struct
func extractAudioStream(data *ffprobe.ProbeData, mediaData *MediaData) error {
	// Loop through all audio streams in the ffprobe data
	for _, audioStream := range data.StreamType(ffprobe.StreamAudio) {
		// Extract the language of the audio stream from the ffprobe data
		language, err := audioStream.TagList.GetString("language")
		if err != nil {
			// If the language tag is not available, try to extract the title tag
			language, err = audioStream.TagList.GetString("title")
			if err != nil {
				// If neither tag is available, set the language to "Unknown"
				language = "Unknown"
			}
		}

		// Parse the bitrate of the audio stream
		bitrate, err := strconv.ParseFloat(audioStream.BitRate, 64)
		if err != nil {
			log.Println("Error parsing audio bitrate:", err)
			bitrate = 0
		}

		// Add a new AudioData struct to the mediaData.audio array
		mediaData.Audios = append(mediaData.Audios, AudioData{
			Codec:    strings.ToUpper(audioStream.CodecName),
			Language: language,
			Bitrate:  bitrate,
		})
	}
	return nil
}
