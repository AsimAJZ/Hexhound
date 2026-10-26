package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
)

type CheckResult struct {
	Name    string
	Passed  bool
	Details string
}

type IHDRreport struct {
	DataSize    uint32
	ChunkType   string
	ChunkData   []byte
	ChunkCrc    uint32
	ComputedCrc uint32
	Width       uint32
	Height      uint32
	Bitdepth    uint8
	ColorType   uint8
	Compression uint8
	Filter      uint8
	Interlace   uint8

	Checks []CheckResult
}

type IDATreport struct {
	DataSize    uint32
	ChunkType   string
	ChunkData   []byte
	ChunkCrc    uint32
	ComputedCrc uint32

	Checks []CheckResult
}

const maxChunkSize uint32 = 16 * 1024 * 1024 // 16 MiB

func IhdrBoundaryCheck(rep *IHDRreport) {

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "dataSize",
		Passed:  rep.DataSize == 13,
		Details: fmt.Sprintf("(dataSize=%d is %s)", rep.DataSize, declarator(rep.DataSize == 13)),
	})

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "chunkType",
		Passed:  rep.ChunkType == "IHDR",
		Details: fmt.Sprintf("(chunkType=%s is %s)", rep.ChunkType, declarator(rep.ChunkType == "IHDR")),
	})

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "chunkCrc",
		Passed:  rep.ChunkCrc == rep.ComputedCrc,
		Details: fmt.Sprintf("(CRC is %s)", declarator(rep.ChunkCrc == rep.ComputedCrc)),
	})

	validWidth := checkDimensions(rep.Width)
	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "width",
		Passed:  validWidth,
		Details: fmt.Sprintf("(Width : %d is %s)", rep.Width, declarator(validWidth)),
	})
	validheight := checkDimensions(rep.Height)
	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "height",
		Passed:  validheight,
		Details: fmt.Sprintf("(height : %d is %s)", rep.Height, declarator(validheight)),
	})

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "compression",
		Passed:  rep.Compression == 0,
		Details: fmt.Sprintf("(compression=%d is %s)", rep.Compression, declarator(rep.Compression == 0)),
	})

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "filter",
		Passed:  rep.Filter == 0,
		Details: fmt.Sprintf("(filter=%d is %s)", rep.Filter, declarator(rep.Filter == 0)),
	})

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "interlace",
		Passed:  rep.Interlace == 0 || rep.Interlace == 1,
		Details: fmt.Sprintf("(interlace=%d is %s)", rep.Interlace, declarator(rep.Interlace == 0 || rep.Interlace == 1)),
	})

	validColorType := rep.ColorType == 0 || rep.ColorType == 2 || rep.ColorType == 3 || rep.ColorType == 4 || rep.ColorType == 6
	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "colorType",
		Passed:  validColorType,
		Details: fmt.Sprintf("(Color type %d is %s)", rep.ColorType, declarator(validColorType)),
	})

	validDepth := validColorType &&
		validBitDepth(rep.ColorType, rep.Bitdepth)
	rep.Checks = append(rep.Checks, CheckResult{
		Name:   "bitDepth",
		Passed: validDepth,
		Details: fmt.Sprintf(
			"Bit depth %d is %s for color type %d",
			rep.Bitdepth,
			map[bool]string{true: "valid", false: "invalid"}[validDepth],
			rep.ColorType,
		),
	})
}

func IhdrCheck(pack *packet) (IHDRreport, error) {
	rep := IHDRreport{}
	if len(pack.Prefix) < 33 {
		return rep, fmt.Errorf("Insufficient bytes found in IHDR field.")
	}

	chunkLen := binary.BigEndian.Uint32(pack.Prefix[8:12])
	chunkType := pack.Prefix[12:16]
	chunkData := pack.Prefix[16:29]

	width := binary.BigEndian.Uint32(pack.Prefix[16:20])
	height := binary.BigEndian.Uint32(pack.Prefix[20:24])
	bitDepth := pack.Prefix[24]
	colorType := pack.Prefix[25]
	compression := pack.Prefix[26]
	fitler := pack.Prefix[27]
	interlace := pack.Prefix[28]

	buff := append(chunkType, chunkData...)
	computedCrc := crc32.Checksum(buff, crc32.IEEETable)
	chunkCrc := binary.BigEndian.Uint32(pack.Prefix[29:33])

	rep = IHDRreport{
		DataSize:    chunkLen,
		ChunkType:   string(chunkType),
		ChunkData:   chunkData,
		Width:       width,
		Height:      height,
		Bitdepth:    uint8(bitDepth),
		ColorType:   colorType,
		Compression: compression,
		Filter:      fitler,
		Interlace:   interlace,
		ChunkCrc:    chunkCrc,
		ComputedCrc: computedCrc,
	}

	IhdrBoundaryCheck(&rep)

	return rep, nil
}

func IdatBoundaryCheck(rep *IDATreport) {

	computedLen := uint32(len(rep.ChunkData))
	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "dataSize",
		Passed:  rep.DataSize == computedLen,
		Details: fmt.Sprintf("dataSize=%d is %s", rep.DataSize, declarator(rep.DataSize == computedLen)),
	})

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "chunkType",
		Passed:  rep.ChunkType == "IDAT",
		Details: fmt.Sprintf("chunkType:%s is %s", rep.ChunkType, declarator(rep.ChunkType == "IDAT")),
	})

	rep.Checks = append(rep.Checks, CheckResult{
		Name:    "chunkCrc",
		Passed:  rep.ChunkCrc == rep.ComputedCrc,
		Details: fmt.Sprintf("CRC is %s", declarator(rep.ChunkCrc == rep.ComputedCrc)),
	})

}

func IDATcheck(pack *packet, IhdrReport *IHDRreport) ([]IDATreport, error) {
	var reports []IDATreport
	var totalIdatChunkData []byte

	file, err := os.Open(pack.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("Opening file for IDAT check failed: %w", err)
	}

	defer file.Close()

	offset, err := file.Seek(8, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("Offsetting failure at %d", offset)
	}

	buffer_len := make([]byte, 4)
	buffer_type := make([]byte, 4)
	buffer_crc := make([]byte, 4)

	seenIDAT := false
	IdatSequenceEnded := false

	for {

		_, err = io.ReadFull(file, buffer_len)
		if err != nil {
			return nil, fmt.Errorf("failed chunktype reading: %w", err)
		}
		length := binary.BigEndian.Uint32(buffer_len)

		currentOffset, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("getting current offset : %w", err)
		}
		remaining := pack.FileSize - currentOffset
		if remaining < 8 || int64(length) > remaining-8 {
			return nil, fmt.Errorf("chunk lenght exceeds remaining file data")
		}

		_, err = io.ReadFull(file, buffer_type)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("failed chunktype reading: %w", err)
		}
		chunkType := string(buffer_type)

		if chunkType == "IDAT" {
			if IdatSequenceEnded {
				return nil, fmt.Errorf("IDAT chunk found after sequence ended.")
			}
			seenIDAT = true
			if length > maxChunkSize {
				return nil, fmt.Errorf("chunk length exceeds HexHound resource policy")
			}
			buffer_data := make([]byte, int(length))
			_, err := io.ReadFull(file, buffer_data)
			if err != nil {
				return nil, err
			}
			_, err = io.ReadFull(file, buffer_crc)
			if err != nil {
				return nil, err
			}
			buff := append(buffer_type, buffer_data...)
			computedCrc := crc32.Checksum(buff, crc32.IEEETable)
			chunkCrc := binary.BigEndian.Uint32(buffer_crc)
			rep := IDATreport{
				DataSize:    length,
				ChunkType:   chunkType,
				ChunkData:   buffer_data,
				ChunkCrc:    chunkCrc,
				ComputedCrc: computedCrc,
			}
			IdatBoundaryCheck(&rep)
			totalIdatChunkData = append(totalIdatChunkData, rep.ChunkData...)
			reports = append(reports, rep)

			continue
		}

		if seenIDAT {
			IdatSequenceEnded = true
		}

		if chunkType == "IEND" {
			if length != 0 {
				return nil, fmt.Errorf("IEND chunk must have zero data length")
			}

			// Skip the IEND CRC.
			if _, err = file.Seek(4, io.SeekCurrent); err != nil {
				return nil, fmt.Errorf("reading past IEND CRC: %w", err)
			}

			if !seenIDAT {
				return nil, fmt.Errorf("PNG contains no IDAT chunks")
			}
			// validated chunk data will be added here to the reports using a loop
			data, err := dataDecompressor(totalIdatChunkData)
			if err != nil {
				return nil, err
			}
			err = filterBytesValidator(data, IhdrReport)
			if err != nil {
				return nil, err
			}
			log.Println("decompressed Data size: ", len(data))
			return reports, nil
		}

		_, err = file.Seek(int64(length)+4, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("skipping chunk %s: %w", chunkType, err)
		}

	}
}

//chunk Data Validation in IHDR needed
