// Package server provides HTTP server for image files.
//
//TODO: use this service as a proxy for other services like imgix, imageflow,
// so that we will have a consistent API. Requests to this server will be
// processed and then the request will be redirected to the actual server
// with the appropriate URL.
//TODO: compatibility mode: translation from parameters for other image servers.
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

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/cespare/xxhash"
	"github.com/gabriel-vasile/mimetype"
	"github.com/richardlehane/crock32"

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
		Width:        widthDefault,
		Height:       heightDefault,
		Fit:          paramsFitDefault,
		Scale:        paramsScaleDefault,
		PaddingColor: nil,
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

	// We use this to identify a variant. It's constructed from a normalized
	// parameter set. This will be appended to the filename.
	variantKey := handler.variantKeyFromParams(params)

	reqPath := r.URL.Path
	processedFilePath := filepath.
		Join(handler.config.RawFilesDir, reqPath+"_1"+variantKey)

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

	srcImageFilePath := filepath.
		Join(handler.config.RawFilesDir, reqPath)

	mime, err := mimetype.DetectFile(srcImageFilePath)
	if err != nil || mime == nil {
		//TODO: check the error
		log.Warn().Err(err).
			Msgf("Unable to detect MIME of the file: %s", srcImageFilePath)
		http.Error(w, http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}

	srcImage, err := imgio.Open(srcImageFilePath)
	if err != nil {
		//TODO: check the error
		log.Warn().Err(err).
			Msgf("Unable to open the file: %s", srcImageFilePath)
		http.Error(w, http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}

	srcImageWidth := srcImage.Bounds().Dx()
	srcImageHeight := srcImage.Bounds().Dy()
	srcImageAspectRatio := float32(srcImageWidth) / float32(srcImageHeight)

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
		if srcImageWidth <= params.Width && srcImageHeight <= params.Height {
			// Pad if requested, otherwise simply save into the new location.
			if params.IsPaddingRequested() {
				panic("TODO")
			} else {
				saveImageAndRespond(srcImage)
				return
			}
		}
		if params.Fit == FitModeContain {
			wRatio := float32(srcImageWidth) / float32(params.Width)
			hRatio := float32(srcImageHeight) / float32(params.Height)
			if params.IsPaddingRequested() {
				panic("TODO: pad")
			} else {
				var scaledWidth, scaledHeight int
				if wRatio > hRatio {
					scaledWidth = params.Width
					scaledHeight = int(float32(scaledWidth) / srcImageAspectRatio)
				} else {
					scaledHeight = params.Height
					scaledWidth = int(float32(scaledHeight) * srcImageAspectRatio)
				}
				outImage := transform.Resize(srcImage,
					scaledWidth, scaledHeight,
					imageScaleDownAlg)
				saveImageAndRespond(outImage)
				return
			}
		} else if params.Fit == FitModeCrop {
			// - simply save if the image is smaller or has the same dimensions
			//   as requested.
			// - don't resize if the non-cropped side size is less than or
			//   equal to the requested size. pad if requested.
			wRatio := float32(srcImageWidth) / float32(params.Width)
			hRatio := float32(srcImageHeight) / float32(params.Height)
			var scaledWidth, scaledHeight int
			if wRatio <= hRatio {
				scaledWidth = params.Width
				scaledHeight = int(float32(scaledWidth) / srcImageAspectRatio)
			} else {
				if srcImageHeight < params.Height {
					scaledHeight = srcImageHeight
				} else {
					scaledHeight = params.Height
				}
				scaledWidth = int(float32(scaledHeight) * srcImageAspectRatio)
			}
			workImage := transform.Resize(srcImage, scaledWidth, scaledHeight,
				imageScaleDownAlg)
			if scaledHeight <= params.Height && scaledWidth <= params.Width {
				if params.IsPaddingRequested() {
					panic("TODO: pad")
				} else {
					saveImageAndRespond(workImage)
					return
				}
			}
			var cropSourceRect image.Rectangle
			if scaledWidth > params.Width {
				wDiff := scaledWidth - params.Width
				offX := wDiff / 2
				cropSourceRect = image.Rectangle{
					Min: image.Point{X: offX, Y: 0},
					Max: image.Point{
						X: params.Width + offX,
						Y: params.Height},
				}
			} else {
				hDiff := scaledHeight - params.Height
				offY := hDiff / 2
				cropSourceRect = image.Rectangle{
					Min: image.Point{X: 0, Y: offY},
					Max: image.Point{
						X: params.Width,
						Y: params.Height + offY},
				}
			}
			workImage = transform.Crop(workImage, cropSourceRect)
			saveImageAndRespond(workImage)
			return
		}
	} else if params.Scale == ScaleDirectionUp {
		if srcImageWidth >= params.Width && srcImageHeight >= params.Height {
			// Process the file by removing any metadata
		}
	}

	http.Error(w, http.StatusText(http.StatusNotImplemented),
		http.StatusNotImplemented)
	return
}

func (handler *Handler) variantKeyFromParams(
	params processingParameters,
) string {
	encodedParams := params.Encode()
	h := xxhash.Sum64String(encodedParams)
	return crock32.Encode(h)
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

	// PaddingColor is used to fill the canvas outside the projected image.
	// If this value is not provided, we won't pad the image. Use
	// IsPaddingRequested method for determining if padding was requested.
	PaddingColor *color.RGBA
}

func (params processingParameters) Encode() string {
	values := url.Values{}
	values.Set("width", strconv.FormatInt(int64(params.Width), 10))
	values.Set("height", strconv.FormatInt(int64(params.Height), 10))
	values.Set("fit", params.Fit.String())
	values.Set("scale", params.Scale.String())
	if params.PaddingColor != nil {
		values.Set("padcol", "#"+rgbaToARGBHex(*params.PaddingColor))
	}
	return values.Encode()
}

// IsPaddingRequested returns true if padding was requested.
func (params processingParameters) IsPaddingRequested() bool {
	return params.PaddingColor != nil && params.PaddingColor.A > 0
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
