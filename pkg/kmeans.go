package fimgs

// #cgo CFLAGS: -g -Wall -march=native -O3 -Wconversion
// #include "fastkmeans.h"
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"unsafe"
)

func makeColorArray(len int) [][3]int64 {
	return make([][3]int64, len)
}

func minkowskiiDist(a, b [3]int64) int64 {
	pa := (*C.long)(unsafe.Pointer(&a[0]))
	pb := (*C.long)(unsafe.Pointer(&b[0]))
	return int64(C.fastDist(pa, pb))
}

func initClusterCenters(pixelColors [][3]int64, clustersCount int) [][3]int64 {
	clustersCenters := makeColorArray(clustersCount)
	clustersCenters[clustersCount-1] = pixelColors[rand.Intn(len(pixelColors))]
	minClusterDistance := make([]int64, len(pixelColors))
	minClusterDistanceSum := int64(0)
	for i, pixelColor := range pixelColors {
		minClusterDistance[i] = minkowskiiDist(pixelColor, clustersCenters[clustersCount-1])
		minClusterDistanceSum += minClusterDistance[i]
	}
	for k := 0; k < clustersCount-1; k++ {
		x := rand.Int63n(minClusterDistanceSum)
		for i, pixelColor := range pixelColors {
			x -= minClusterDistance[i]
			if x < 0 {
				clustersCenters[k] = pixelColor
				break
			}
		}
		if k == clustersCount-1 {
			break
		}
		for i, pixelColor := range pixelColors {
			if newDistance := minkowskiiDist(pixelColor, clustersCenters[0]); newDistance < minClusterDistance[i] {
				minClusterDistanceSum += newDistance - minClusterDistance[i]
				minClusterDistance[i] = newDistance
			}
		}
	}
	return clustersCenters
}

func kmeansIters(clustersCenters, pixelColors [][3]int64, clustersCount int) {
	C.kmeansIters(
		(*[3]C.long)(unsafe.Pointer(&clustersCenters[0][0])),
		C.int(clustersCount),
		(*[3]C.long)(unsafe.Pointer(&pixelColors[0][0])),
		C.int(len(pixelColors)),
	)
}

func ApplyKMeans(im image.Image, clustersCount int) image.RGBA {
	// TODO: use one or another kMeansIters depending on having AVX2
	// fmt.Println("has AVX2: ", cpuid.CPU.Supports(cpuid.AVX2))
	imageWidth := im.Bounds().Dx()
	pixelColors := makeColorArray(imageWidth * im.Bounds().Dy())
	for j := 0; j < im.Bounds().Dy(); j++ {
		for i := 0; i < im.Bounds().Dx(); i++ {
			k := i + j*imageWidth
			r, g, b, _ := im.At(i, j).RGBA()
			pixelColors[k] = [3]int64{int64(r), int64(g), int64(b)}
		}
	}
	rand.Seed(0)
	clustersCenters := initClusterCenters(pixelColors, clustersCount)
	// TODO: try to sample mini-batches (random subdatasets)
	kmeansIters(clustersCenters, pixelColors, clustersCount)
	filtered_im := image.NewRGBA(im.Bounds())
	for j := 0; j < im.Bounds().Dy(); j++ {
		for i := 0; i < im.Bounds().Dx(); i++ {
			pixel := pixelColors[i+j*imageWidth]
			minCluster := 0
			minDist := minkowskiiDist(pixel, clustersCenters[0])
			for k := 1; k < clustersCount; k++ {
				dist := minkowskiiDist(pixel, clustersCenters[k])
				if dist < minDist {
					minCluster = k
					minDist = dist
				}
			}
			filtered_im.Set(i, j, color.RGBA{
				uint8((clustersCenters[minCluster][0] & 0xFF00) >> 8),
				uint8((clustersCenters[minCluster][1] & 0xFF00) >> 8),
				uint8((clustersCenters[minCluster][2] & 0xFF00) >> 8),
				255,
			})
		}
	}
	return *filtered_im
}

// TODO: filter init is also validation?
func ApplyKMeansFilter(sourceImageFilename string, resultImageFilename string, clustersCount int) error {
	if clustersCount < 2 {
		return fmt.Errorf("'n' must be at least 2, you gave n=%d", clustersCount)
	}

	im, err := LoadImageFile(sourceImageFilename)
	if err != nil {
		return fmt.Errorf("error occured while loading image: %q", err)
	}

	filtered_im := ApplyKMeans(im, clustersCount)

	if err := saveImage(filtered_im, resultImageFilename); err != nil {
		return err
	}

	return nil
}
