package fimgs

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
)

func makeColorArray(len int) [][]int64 {
	data := make([]int64, len*3)
	res := make([][]int64, len)
	for i := 0; i < len; i++ {
		res[i] = data[i*3 : i*3+3]
	}
	return res
}

func abs64(x int64) int64 {
	mask := x >> 63
	return (x + mask) ^ mask
}

func minkowskiiDist(a, b []int64) int64 {
	return abs64(a[0]-b[0]) + abs64(a[1]-b[1]) + abs64(a[2]-b[2])
}

func initClusterCenters(pixelColors [][]int64, clustersCount int) [][]int64 {
	clustersCenters := makeColorArray(clustersCount)
	copy(clustersCenters[0], pixelColors[rand.Intn(len(pixelColors))])
	minClusterDistance := make([]int64, len(pixelColors))
	minClusterDistanceSum := int64(0)
	for i, pixelColor := range pixelColors {
		minClusterDistance[i] = minkowskiiDist(pixelColor, clustersCenters[0])
		minClusterDistanceSum += minClusterDistance[i]
	}
	for k := 1; k < clustersCount; k++ {
		x := rand.Int63n(minClusterDistanceSum)
		for i, pixelColor := range pixelColors {
			x -= minClusterDistance[i]
			if x < 0 {
				copy(clustersCenters[k], pixelColor)
				break
			}
		}
		if k == clustersCount-1 {
			break
		}
		for i, pixelColor := range pixelColors {
			newDistance := minkowskiiDist(pixelColor, clustersCenters[0])
			if newDistance < minClusterDistance[i] {
				minClusterDistanceSum += newDistance - minClusterDistance[i]
				minClusterDistance[i] = newDistance
			}
		}
	}
	return clustersCenters
}

func kmeansIters(clustersCenters, pixelColors [][]int64, clustersCount int) {
	for epoch := 0; epoch < 300; epoch++ {
		sumAndCount := make([]int64, clustersCount*4) // sum of Rs, Gs, Bs and count
		for _, pixelColor := range pixelColors {
			minCluster := 0
			minDist := minkowskiiDist(pixelColor, clustersCenters[0])
			for k := 1; k < clustersCount; k++ {
				newDist := minkowskiiDist(pixelColor, clustersCenters[k])
				if newDist < minDist {
					minCluster = k
					minDist = newDist
				}
			}
			sumAndCount[minCluster*4+0] += pixelColor[0]
			sumAndCount[minCluster*4+1] += pixelColor[1]
			sumAndCount[minCluster*4+2] += pixelColor[2]
			sumAndCount[minCluster*4+3]++
		}
		movement := int64(0)
		for i := 0; i < clustersCount; i++ {
			count := sumAndCount[i*4+3]
			if count == 0 {
				continue
			}
			sumAndCount[i*4+0] /= count
			sumAndCount[i*4+1] /= count
			sumAndCount[i*4+2] /= count
			movement += minkowskiiDist(clustersCenters[i], sumAndCount[i*4:i*4+4])
			clustersCenters[i] = sumAndCount[i*4 : i*4+4]
		}
		if movement < 100 {
			break
		}
	}
}

func ApplyKMeans(im image.Image, clustersCount int) image.Image {
	imageWidth := im.Bounds().Dx()
	pixelColors := makeColorArray(imageWidth * im.Bounds().Dy())
	for i := im.Bounds().Min.X; i < im.Bounds().Max.X; i++ {
		for j := im.Bounds().Min.Y; j < im.Bounds().Max.Y; j++ {
			k := i + j*imageWidth
			r, g, b, _ := im.At(i, j).RGBA()
			r64, g64, b64 := int64(r), int64(g), int64(b)
			pixelColors[k][0] = r64
			pixelColors[k][1] = g64
			pixelColors[k][2] = b64
		}
	}
	rand.Seed(0)
	clustersCenters := initClusterCenters(pixelColors, clustersCount)
	// TODO: try to sample mini-batches (random subdatasets)
	kmeansIters(clustersCenters, pixelColors, clustersCount)
	filtered_im := image.NewRGBA(im.Bounds())
	for i := im.Bounds().Min.X; i < im.Bounds().Max.X; i++ {
		for j := im.Bounds().Min.Y; j < im.Bounds().Max.Y; j++ {
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
				uint8((clustersCenters[minCluster][0]) / 0x100),
				uint8((clustersCenters[minCluster][1]) / 0x100),
				uint8((clustersCenters[minCluster][2]) / 0x100),
				255,
			})
		}
	}
	return filtered_im
}

// TODO: filter init is also validation?
func ApplyKMeansFilter(sourceImageFilename string, resultImageFilename string, clustersCount int) (err error) {
	// f, _ := os.Create("cpu.pb")
	// defer f.Close() // error handling omitted for example
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	if clustersCount < 2 {
		return fmt.Errorf("'n' must be at least 2, you gave n=%d", clustersCount)
	}
	im, err := LoadImageFile(sourceImageFilename)
	if err != nil {
		return fmt.Errorf("error occured while loading image:\n%q", err)
	}
	filtered_im := ApplyKMeans(im, clustersCount)
	err = saveImage(filtered_im, resultImageFilename)
	if err != nil {
		return
	}
	return
}
