package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"math"
)

const maxDataSize uint32 = 64 * 1024 * 1024 // 64 MiB

func validBitDepth(colorType, bitDepth uint8) bool {
	switch colorType {
	case 0:
		return bitDepth == 1 || bitDepth == 2 ||
			bitDepth == 4 || bitDepth == 8 || bitDepth == 16
	case 2:
		return bitDepth == 8 || bitDepth == 16
	case 3:
		return bitDepth == 1 || bitDepth == 2 ||
			bitDepth == 4 || bitDepth == 8
	case 4, 6:
		return bitDepth == 8 || bitDepth == 16
	default:
		return false
	}
}

func checkDimensions(dim uint32) bool {
	if dim >= 1 && dim <= 2147483647 {
		return true
	}
	return false
}

func declarator(b bool) string {
	if b {
		return "valid"
	}
	return "invalid"
}

func dataDecompressor(data []byte) ([]byte, error) {
	compressedData := bytes.NewReader(data)

	zlibRead, err := zlib.NewReader(compressedData)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize zlib Reader for IDAT chunk Data: %w", err)
	}

	defer zlibRead.Close()

	var decompressedData bytes.Buffer

	limitedReader := io.LimitReader(zlibRead, int64(maxDataSize)+1)

	_, err = io.Copy(&decompressedData, limitedReader)
	if err != nil {
		return nil, fmt.Errorf("Decompression failed: %w", err)
	}

	rawData := decompressedData.Bytes()

	if len(rawData) > int(maxDataSize) {
		return nil, fmt.Errorf("IDAT chunks data size exceeds HexHound resources policy.")
	}

	return rawData, nil
}

func filterBytesValidator(data []byte, IhdrRep *IHDRreport) error {
	var channels int
	interlace := IhdrRep.Interlace
	colorType := IhdrRep.ColorType
	bitDepth := int(IhdrRep.Bitdepth)
	width := int(IhdrRep.Width)
	height := int(IhdrRep.Height)

	switch interlace {
	case 0:
		switch colorType {
		case 0:
			channels = 1
		case 2:
			channels = 3
		case 3:
			channels = 1
		case 4:
			channels = 2
		case 6:
			channels = 4
		default:
			return fmt.Errorf("Unsupported ColorType of png file found: %d", colorType)
		}
		if bitDepth == 0 {
			return fmt.Errorf("Invalid bitdepth value, reported in IDAT filterBytesValidator.")
		}
		bitsPerPixel := channels * bitDepth

		if width > (math.MaxInt-7)/bitsPerPixel {
			return fmt.Errorf("Exceeded Resource limit of HexHound policy: filterByteValidator: width = %d", width)
		}
		rowBytes := (width*bitsPerPixel + 7) / 8
		scanlineSize := rowBytes + 1
		if height == 0 {
			return fmt.Errorf("Invalid height value, found in IDAT filterBytesValidator.")
		}
		if scanlineSize > int(maxDataSize)/height {
			return fmt.Errorf("Exceeded Resource limit of HexHound policy: filterByteValidator: scanlineSize = %d, height = %d", scanlineSize, height)
		}
		expectedSize := scanlineSize * height
		if expectedSize != len(data) {
			return fmt.Errorf("Decompressed data has an invalid Size.")
		}
		scanlineNum := 0
		for i := 0; i < expectedSize; i += scanlineSize {
			if data[i] > 4 {
				return fmt.Errorf("Invalid filter-byte encountered at postion %d at scanLin %d with value: %d", i, scanlineNum, data[i])
			}
			scanlineNum++
		}
	case 1:
		return fmt.Errorf("This type of png is unsupported for IDAT validation")
	default:
		return fmt.Errorf("Invalid interlace value, found in IDAT filterBytesValidator.")
	}
	return nil
}
