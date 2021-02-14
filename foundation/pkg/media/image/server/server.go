// Package server provides HTTP server for image files.
//
//TODO: use this service as a proxy for other services like imgix, imageflow,
// so that we will have a consistent API. Requests to this server will be
// processed and then the request will be redirected to the actual server
// with the appropriate URL.
package server

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gabriel-vasile/mimetype"
	"github.com/richardlehane/crock32"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/errors"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/logging"
)

var log = logging.NewPkgLogger()

const (
	paramsHeightDefault = 1024
	paramsWidthDefault  = 1024
	paramsFitDefault    = FitModeContain
	paramsScaleDefault  = ScaleDirectionDown
)

type HandlerConfig struct {
	ProcessedFilesDir string `env:"PROCESSED_FILES_DIR`
	RawFilesDir       string `env:"-"`
	HeightDefault     int32  `env:"HEIGHT_DEFAULT"`
	WidthDefault      int32  `env:"WIDTH_DEFAULT"`
}

func NewHandler(config HandlerConfig) (*Handler, error) {
	return &Handler{config: config}, nil
}

type Handler struct {
	config HandlerConfig
}

// ServeHTTP conforms Go's HTTP Handler interface.
func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	widthDefault := paramsWidthDefault
	if handler.config.WidthDefault > 0 {
		widthDefault = int(handler.config.WidthDefault)
	}
	heightDefault := paramsWidthDefault
	if handler.config.HeightDefault > 0 {
		heightDefault = int(handler.config.HeightDefault)
	}

	params := processingParameters{
		Width:    widthDefault,
		Height:   heightDefault,
		Fit:      paramsFitDefault,
		Scale:    paramsScaleDefault,
		PadColor: nil,
	}

	reqQuery := r.URL.Query()
	for queryKey, queryValues := range reqQuery {
		queryValue := queryValues[0] //TODO: might be empty
		switch strings.ToLower(queryKey) {
		case "width", "w":
			wv, err := strconv.ParseInt(queryValue, 10, 32)
			if err != nil {
				log.Warn().Err(err).
					Msgf("Query width cannot be parsed. Value is %q", queryValue)
				http.Error(w, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				return
			}
			params.Width = int(wv)
		case "height", "h":
			hv, err := strconv.ParseInt(queryValue, 10, 32)
			if err != nil {
				log.Warn().Err(err).
					Msgf("Query height cannot be parsed. Value is %q", queryValue)
				http.Error(w, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				return
			}
			params.Height = int(hv)
		case "fit":
			fit, err := FitModeFromString(queryValues[0])
			if err != nil {
				log.Warn().Err(err).
					Msgf("Query fit cannot be parsed. Value is %q", queryValue)
				http.Error(w, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				return
			}
			if fit != FitModeUnspecified {
				params.Fit = fit
			}
		case "scale":
			scale, err := ScaleDirectionFromString(queryValues[0])
			if err != nil {
				log.Warn().Err(err).
					Msgf("Query scale cannot be parsed. Value is %q", queryValue)
				http.Error(w, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				return
			}
			if scale != ScaleDirectionUnspecified {
				params.Scale = scale
			}
		default:
			log.Warn().
				Msgf("Unsupported query key %q", queryKey)
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
	}

	//TODO: validate params here
	if params.Width <= 0 {
		if params.Width < 0 {
			log.Warn().
				Msgf("Invalid width value %v", params.Width)
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
		params.Width = widthDefault
	}
	if params.Height <= 0 {
		if params.Height < 0 {
			log.Warn().
				Msgf("Invalid height value: %v", params.Height)
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
		params.Height = heightDefault
	}
	if params.Fit == FitModeUnspecified {
		params.Fit = paramsFitDefault
	}
	if params.Scale == ScaleDirectionUnspecified {
		params.Scale = paramsScaleDefault
	}

	encodedParams := params.Encode()
	h := xxhash.ChecksumString64(encodedParams)
	paramsKey := crock32.Encode(h)

	reqPath := r.URL.Path
	processedFilePath := filepath.
		Join(handler.config.RawFilesDir, reqPath+"_1"+paramsKey)
	processedFile, err := os.Open(processedFilePath)
	if err != nil {
		pathErr, ok := err.(*os.PathError)
		if !ok || pathErr == nil {
			panic(err)
		}
		if !errors.Is(pathErr.Err, os.ErrNotExist) {
			fmt.Printf("%#v\n", pathErr.Err)
			panic(pathErr)
		}
	}
	if processedFile != nil {
		processedFile.Close()
		http.ServeFile(w, r, processedFilePath)
		return
	}

	rawFilePath := filepath.
		Join(handler.config.RawFilesDir, reqPath)

	mime, err := mimetype.DetectFile(rawFilePath)
	if err != nil || mime == nil {
		//TODO: check the error
		log.Warn().Err(err).
			Msgf("Unable to detect MIME of the file: %s", rawFilePath)
		http.Error(w, http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}

	rawImage, err := imgio.Open(rawFilePath)
	if err != nil {
		//TODO: check the error
		log.Warn().Err(err).
			Msgf("Unable to open the file: %s", rawFilePath)
		http.Error(w, http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}

	rawImageW := rawImage.Bounds().Dx()
	rawImageH := rawImage.Bounds().Dy()
	rawImageAR := float32(rawImageW) / float32(rawImageH)

	//TODO: support gif (animated?) and webp, and optimize the files.
	saveImageAndRespond := func(img image.Image) {
		var encoder imgio.Encoder
		if mime.Is("image/png") || mime.Is("image/bmp") {
			encoder = imgio.PNGEncoder()
		} else if mime.Is("image/jpeg") {
			encoder = imgio.JPEGEncoder(89)
		}
		if encoder == nil {
			panic("Unsupported image type")
		}
		err = imgio.Save(processedFilePath, img, encoder)
		if err != nil {
			panic("Unable to save the image")
		}
		http.ServeFile(w, r, processedFilePath)
	}

	if params.Scale == ScaleDirectionDown {
		if rawImageW <= params.Width && rawImageH <= params.Height {
			//TODO: check if padding we want some padding, otherwise
			// simply save into the new location.
			if params.PadColor != nil {
				panic("TODO")
			} else {
				saveImageAndRespond(rawImage)
				return
			}
		}
		if params.Fit == FitModeContain {
			wRatio := float32(rawImageW) / float32(params.Width)
			hRatio := float32(rawImageH) / float32(params.Height)
			//TODO: Pad
			if params.PadColor != nil {
			} else {
				var targetWidth, targetHeight int
				if wRatio > hRatio {
					targetWidth = params.Width
					targetHeight = int(float32(targetWidth) / rawImageAR)
				} else {
					targetHeight = params.Height
					targetWidth = int(float32(targetHeight) * rawImageAR)
				}
				outImage := transform.Resize(rawImage,
					targetWidth, targetHeight,
					imageScaleDownAlg)
				saveImageAndRespond(outImage)
				return
			}
		} else if params.Fit == FitModeCrop {
			// pad is not applicable
			wRatio := float32(rawImageW) / float32(params.Width)
			hRatio := float32(rawImageH) / float32(params.Height)
			var targetWidth, targetHeight int
			if wRatio <= hRatio {
				targetWidth = params.Width
				targetHeight = int(float32(targetWidth) / rawImageAR)
			} else {
				targetHeight = params.Height
				targetWidth = int(float32(targetHeight) * rawImageAR)
			}
			tmpImage := transform.Resize(rawImage, targetWidth, targetHeight,
				imageScaleDownAlg)
			if targetHeight == params.Height {
				if targetWidth == params.Width {
					// No need to crop
					saveImageAndRespond(tmpImage)
					return
				}
				wDiff := targetWidth - params.Width
				offX := wDiff / 2
				tmpImage = transform.Crop(tmpImage, image.Rectangle{
					Min: image.Point{X: offX, Y: 0},
					Max: image.Point{
						X: params.Width + offX,
						Y: params.Height},
				})
				saveImageAndRespond(tmpImage)
				return
			}
			hDiff := targetHeight - params.Height
			offY := hDiff / 2
			tmpImage = transform.Crop(tmpImage, image.Rectangle{
				Min: image.Point{X: 0, Y: offY},
				Max: image.Point{
					X: params.Width,
					Y: params.Height + offY},
			})
			saveImageAndRespond(tmpImage)
			return
		}
	} else if params.Scale == ScaleDirectionUp {
		if rawImageW >= params.Width && rawImageH >= params.Height {
			// Process the file by removing any metadata
		}
	}

	http.Error(w, http.StatusText(http.StatusNotImplemented),
		http.StatusNotImplemented)
	return

}

var (
	imageScaleDownAlg = transform.Lanczos
	imageScaleUpAlg   = transform.MitchellNetravali
)

type processingParameters struct {
	Width  int
	Height int
	Fit    FitMode
	Scale  ScaleDirection

	// PadColor is used to fill the canvas outside the projected image.
	// If this value is not provided, we won't pad the image.
	PadColor *color.RGBA
}

func (params processingParameters) Encode() string {
	values := url.Values{}
	values.Set("width", strconv.FormatInt(int64(params.Width), 10))
	values.Set("height", strconv.FormatInt(int64(params.Height), 10))
	values.Set("fit", params.Fit.String())
	values.Set("scale", params.Scale.String())
	if params.PadColor != nil {
		values.Set("pad", "#"+rgbaToARGBHex(*params.PadColor))
	}
	return values.Encode()
}

type FitMode int

const (
	FitModeUnspecified FitMode = iota
	FitModeContain
	FitModeCrop
)

func FitModeFromString(s string) (FitMode, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "contain", "max":
		return FitModeContain, nil
	case "crop", "cover", "min":
		return FitModeCrop, nil
	case "":
		return FitModeUnspecified, nil
	}
	return FitModeUnspecified, errors.Msg("unsupported string value")
}

func (fitMode FitMode) String() string {
	switch fitMode {
	case FitModeContain:
		return "contain"
	case FitModeCrop:
		return "crop"
	case FitModeUnspecified:
		return ""
	}
	return "<invalid>"
}

type ScaleDirection int

const (
	// ScaleDirectionUnspecified is the zero value of ScaleDirection.
	ScaleDirectionUnspecified ScaleDirection = iota
	// ScaleDirectionNone will never scale the image.
	ScaleDirectionNone
	// ScaleDirectionBoth will scale the image up or down to make the
	// image fit into target canvas.
	ScaleDirectionBoth
	// ScaleDirectionUp will only scale up. Images which are smaller than
	// the target canvas will be scaled up while images larger than the
	// canvas will be cropped up.
	ScaleDirectionUp
	// ScaleDirectionDown will only scale down. Images which are larger
	// than the target canvas will be scaled down while images smaller than
	// the canvas will be left as-is.
	ScaleDirectionDown
)

func ScaleDirectionFromString(s string) (ScaleDirection, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "none":
		return ScaleDirectionNone, nil
	case "both":
		return ScaleDirectionBoth, nil
	case "up":
		return ScaleDirectionUp, nil
	case "down":
		return ScaleDirectionDown, nil
	case "":
		return ScaleDirectionUnspecified, nil
	}
	return ScaleDirectionUnspecified, errors.Msg("unsupported string value")
}

func (scaleDir ScaleDirection) String() string {
	switch scaleDir {
	case ScaleDirectionNone:
		return "none"
	case ScaleDirectionBoth:
		return "both"
	case ScaleDirectionUp:
		return "up"
	case ScaleDirectionDown:
		return "down"
	case ScaleDirectionUnspecified:
		return ""
	}
	return "<invalid>"
}

func rgbaToARGBHex(c color.RGBA) string {
	var i uint32
	i |= (uint32(c.A) << 24)
	i |= (uint32(c.R) << 16)
	i |= (uint32(c.G) << 8)
	i |= uint32(c.B)
	return strconv.FormatUint(uint64(i), 16)
}
