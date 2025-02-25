package apis

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/nanoteck137/nosepass"
	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const hlsSegmentLength float64 = 5

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallHandlers(app, g)

	g = router.Group("")
	g.Register(
		pyrin.NormalHandler{
			Name:   "GetPlaylist",
			Method: http.MethodGet,
			Path:   "/:index/index.m3u8",
			HandlerFunc: func(c pyrin.Context) error {
				index, err := strconv.Atoi(c.Param("index"))
				if err != nil {
					return err
				}

				p := "/Users/nanoteck137/anime/Alya Sometimes Hides Her Feelings in Russian S01E01.mkv"
				fmt.Printf("p: %v\n", p)

				ctx := context.TODO()
				data, err := ffprobe.ProbeURL(ctx, p)
				if err != nil {
					return err
				}

				c.Response().Header().Set("Content-Type", "application/vnd.apple.mpegurl")

				duration := data.Format.DurationSeconds

				w := c.Response()

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

					u := ConvertURL(c, fmt.Sprintf("/%d/segment%d.ts", index, segmentIndex))
					fmt.Fprint(w, u + "\n")

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
			Path:   "/:index/:segment",
			HandlerFunc: func(c pyrin.Context) error {
				index, err := strconv.Atoi(c.Param("index"))
				if err != nil {
					return err
				}

				_ = index

				segment := c.Param("segment")
				fmt.Printf("segment: %v\n", segment)

				var segmentIndex int
				_, err = fmt.Sscanf(segment, "segment%d.ts", &segmentIndex)
				if err != nil {
					return err
				}

				// p := "/Users/nanoteck137/anime/Alya Sometimes Hides Her Feelings in Russian S01E01.mkv"
				p := "/Users/nanoteck137/anime/Arifureta Shokugyou de Sekai Saikyou Episode 1.mp4"

				// ctx := context.TODO()
				// data, err := ffprobe.ProbeURL(ctx, p)
				// if err != nil {
				// 	return err
				// }
				//
				// _ = data

				startTime := float64(segmentIndex) * hlsSegmentLength

				sub := "/Users/nanoteck137/projects/transboder/work/metadata/c41e2c95fdd378d4196631ef6330aaae0524b8f7/sub/1.ass"
				_ = sub

				// NOTE(patrik): Based on the code from: 
				// https://github.com/shimberger/gohls/blob/master/internal/hls/segment.go
				args := []string{
					"-nostats",
					"-hide_banner",
					"-loglevel", "warning",
					"-timelimit", "45",

					"-ss", fmt.Sprintf("%v.00", startTime),

					"-hwaccel", "auto",

					"-i", p,

					// "-start_at_zero",
					"-copyts",
					// "-muxdelay", "0",

					"-map", "0:V:0",

					"-t", fmt.Sprintf("%v.00", hlsSegmentLength),

					"-strict", "-2",

					"-ss", fmt.Sprintf("%v.00", startTime),
					// "-filter:v", fmt.Sprintf("subtitles=%s", sub),

					// "-async", "1",

					// 720p
					// "-vf", fmt.Sprintf("scale=-2:%v", res),

					"-vcodec", "libx264",
					"-preset", "veryfast",

					"-map", "0:a:0",

					"-c:a", "aac",
					"-b:a", "128k",
					"-ac", "2",

					// "-pix_fmt", "yuv420p",

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
