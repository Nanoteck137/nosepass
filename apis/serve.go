package apis

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/nanoteck137/nosepass"
	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
)

const hlsSegmentLength float64 = 5

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallSystemHandlers(app, g)
	InstallAuthHandlers(app, g)
	InstallUserHandlers(app, g)
	InstallMediaHandlers(app, g)
	InstallLibraryHandlers(app, g)

	g = router.Group("/api/media")
	g.Register(
		pyrin.NormalHandler{
			Name:   "GetPlaylist",
			Method: http.MethodGet,
			Path:   "/:id/index.m3u8",
			HandlerFunc: func(c pyrin.Context) error {
				id := c.Param("id")

				url := c.Request().URL
				audio, _ := strconv.ParseInt(url.Query().Get("audio"), 10, 0)

				var subtitle int64 = -1

				s := url.Query().Get("subtitle")
				if s != "" {
					subtitle, _ = strconv.ParseInt(s, 10, 0)
				}

				ctx := context.TODO()

				media, err := app.DB().GetMediaById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
					return err
				}

				videoTracks := media.VideoTracks.GetOrEmpty()
				duration := videoTracks[0].Duration

				w := c.Response()
				w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

				// NOTE(patrik): Based on the code from:
				// https://github.com/shimberger/gohls/blob/master/internal/hls/playlist.go
				fmt.Fprint(w, "#EXTM3U\n")
				fmt.Fprint(w, "#EXT-X-VERSION:3\n")
				fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
				fmt.Fprint(w, "#EXT-X-ALLOW-CACHE:YES\n")
				fmt.Fprint(w, "#EXT-X-TARGETDURATION:"+fmt.Sprintf("%.f", hlsSegmentLength)+"\n")
				fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

				leftover := duration
				segmentIndex := 0

				for leftover > 0 {
					if leftover > hlsSegmentLength {
						fmt.Fprintf(w, "#EXTINF: %f,\n", hlsSegmentLength)
					} else {
						fmt.Fprintf(w, "#EXTINF: %f,\n", leftover)
					}

					// fmt.Fprintf(w, getUrl(segmentIndex)+"\n")
					// fmt.Fprintf(w, "segment%d.ts\n", segmentIndex)

					u := ConvertURL(c, fmt.Sprintf("/%v/segment%d.ts?audio=%d&subtitle=%d", id, segmentIndex, audio, subtitle))
					fmt.Fprint(w, u+"\n")

					leftover = leftover - hlsSegmentLength
					segmentIndex++
				}

				fmt.Fprint(w, "#EXT-X-ENDLIST\n")

				return nil
			},
		},
		pyrin.NormalHandler{
			Name:   "GetSegment",
			Method: http.MethodGet,
			Path:   "/:id/:segment",
			HandlerFunc: func(c pyrin.Context) error {
				id := c.Param("id")

				url := c.Request().URL
				audio, _ := strconv.ParseInt(url.Query().Get("audio"), 10, 0)
				subtitle, _ := strconv.ParseInt(url.Query().Get("subtitle"), 10, 0)

				ctx := context.TODO()

				media, err := app.DB().GetMediaById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
					return err
				}

				segment := c.Param("segment")

				var segmentIndex int
				_, err = fmt.Sscanf(segment, "segment%d.ts", &segmentIndex)
				if err != nil {
					return err
				}

				startTime := float64(segmentIndex) * hlsSegmentLength

				mediaDir := app.WorkDir().MediaIdDir(media.Id)

				const videoFormat = "format=yuv420p"

				vfilter := videoFormat
				if subtitle != -1 {
					subtitle := media.Subtitles.GetOrEmpty()[subtitle]
					sub := path.Join(mediaDir.Subtitles(), subtitle.Filename)

					vfilter = fmt.Sprintf("%s,subtitles=%s", videoFormat, sub)
				}

				// NOTE(patrik): Based on the code from:
				// https://github.com/shimberger/gohls/blob/master/internal/hls/segment.go
				args := []string{
					"-nostats",
					"-hide_banner",
					"-loglevel", "warning",
					"-timelimit", "45",

					"-ss", fmt.Sprintf("%v.00", startTime),

					"-hwaccel", "auto",

					"-i", media.Path,

					// "-start_at_zero",
					"-copyts",
					"-muxdelay", "0",

					"-map", "0:V:0",

					"-t", fmt.Sprintf("%v.00", hlsSegmentLength),

					"-strict", "-2",

					"-ss", fmt.Sprintf("%v.00", startTime),
					// "-filter:v", fmt.Sprintf("format=yuv420p,subtitles=%s", sub),
					// "-filter:v", "format=yuv420p",
					"-filter:v", vfilter,

					"-c:v", "libx264",
					"-crf", "18",
					"-preset", "veryfast",
					// "-preset", "ultrafast",

					"-map", fmt.Sprintf("0:a:%d", audio),

					"-c:a", "aac",
					"-b:a", "128k",
					"-ac", "2",

					"-force_key_frames", "expr:gte(t,n_forced*5.000)",

					"-f", "ssegment",
					"-segment_time", fmt.Sprintf("%v.00", hlsSegmentLength),
					// "-initial_offset", fmt.Sprintf("%v.00", startTime),
					"-output_ts_offset", fmt.Sprintf("%v.00", startTime),

					"pipe:out%03d.ts",
				}

				c.Response().Header().Set("Content-Type", "video/mp2t")

				// buffer := bytes.Buffer{}
				cmd := exec.CommandContext(c.Request().Context(), "ffmpeg", args...)
				cmd.Stdout = c.Response()
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					return nil
				}

				return nil
			},
		},
	)
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: nosepass.AppName,
		RegisterHandlers: func(router pyrin.Router) {
			RegisterHandlers(app, router)
		},
	})

	return s, nil
}
